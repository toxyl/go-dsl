package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/chzyer/readline"
	"github.com/toxyl/flo"
)

type paramCompleter struct {
	prefixCompleter readline.PrefixCompleterInterface
	funcParams      map[string][]string // Maps function names to their parameter lists
}

func (c *paramCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Get the original completions
	completions, length := c.prefixCompleter.Do(line, pos)

	// Convert line to string for easier manipulation
	lineStr := string(line)

	// Check if we're inside a function call
	if strings.Contains(lineStr, "(") {
		// Extract the function name
		funcName := strings.Split(lineStr, "(")[0]

		// Get the parameters for this function
		params, exists := c.funcParams[funcName]
		if exists {
			// Count how many parameters have been used
			usedParams := 0
			for _, param := range params {
				if strings.Contains(lineStr, param+"=") {
					usedParams++
				}
			}

			// If all parameters are used, add closing parenthesis to completions
			if usedParams == len(params) {
				// Filter out parameter completions since all are used
				var filteredCompletions [][]rune
				for _, comp := range completions {
					compStr := string(comp)
					// Only keep completions that are closing parenthesis
					if compStr == ")" {
						filteredCompletions = append(filteredCompletions, comp)
					}
				}

				// If no closing parenthesis completion exists, add one
				if len(filteredCompletions) == 0 {
					filteredCompletions = append(filteredCompletions, []rune(")"))
				}

				return filteredCompletions, length
			} else {
				// If not all parameters are used, filter out used parameters
				var filteredCompletions [][]rune
				for _, comp := range completions {
					compStr := string(comp)
					// Keep parameter completions that haven't been used yet
					keep := true
					for _, param := range params {
						if strings.Contains(lineStr, param+"=") && strings.HasPrefix(compStr, param+"=") {
							keep = false
							break
						}
					}
					if keep {
						// Always strip trailing spaces from completions
						filteredCompletions = append(filteredCompletions, []rune(strings.TrimRight(compStr, " ")))
					}
				}
				return filteredCompletions, length
			}
		}
	}

	// For non-function completions, also strip trailing spaces
	var filteredCompletions [][]rune
	for _, comp := range completions {
		filteredCompletions = append(filteredCompletions, []rune(strings.TrimRight(string(comp), " ")))
	}
	return filteredCompletions, length
}

func (c *paramCompleter) GetName() []rune {
	return c.prefixCompleter.GetName()
}

func (c *paramCompleter) GetChildren() []readline.PrefixCompleterInterface {
	return c.prefixCompleter.GetChildren()
}

func (c *paramCompleter) SetChildren(children []readline.PrefixCompleterInterface) {
	c.prefixCompleter.SetChildren(children)
}

func (dsl *dslCollection) shell() {
	debugMode := false

	// Create template data
	type templateData struct {
		Name      string
		Version   string
		Variables []struct {
			Name        string
			Description string
		}
		Functions []struct {
			Name        string
			Description string
			Params      []struct {
				Name    string
				Type    string
				Default any
			}
		}
	}

	data := templateData{
		Name:    dsl.name,
		Version: dsl.version,
	}

	// Add variables
	for name, v := range dsl.vars.data {
		data.Variables = append(data.Variables, struct {
			Name        string
			Description string
		}{
			Name:        name,
			Description: v.meta.desc,
		})
	}

	// Add functions
	for name, fn := range dsl.funcs.data {
		funcData := struct {
			Name        string
			Description string
			Params      []struct {
				Name    string
				Type    string
				Default any
			}
		}{
			Name:        name,
			Description: fn.meta.desc,
		}

		for _, param := range fn.meta.params {
			funcData.Params = append(funcData.Params, struct {
				Name    string
				Type    string
				Default any
			}{
				Name:    param.name,
				Type:    param.typ,
				Default: param.def,
			})
		}

		data.Functions = append(data.Functions, funcData)
	}

	// Parse and execute template
	tmpl, err := template.ParseFS(dslTemplates, "template_shell.tmpl")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(dsl.renderMarkdownToTerminal(buf.String()))
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}
	histFile := filepath.Join(homeDir, fmt.Sprintf("%s_history", dsl.id))

	// Create the base prefix completer
	prefixCompleter := readline.NewPrefixCompleter()

	// Add variable names to completer
	varNames := make([]string, 0, len(dsl.vars.data))
	for name := range dsl.vars.data {
		varNames = append(varNames, name)
	}
	sort.Strings(varNames)
	for _, name := range varNames {
		prefixCompleter.Children = append(prefixCompleter.Children, readline.PcItem(name))
	}

	// Add function names to completer with parameter completion
	funcParams := make(map[string][]string)
	funcNames := make([]string, 0, len(dsl.funcs.data))
	for name := range dsl.funcs.data {
		funcNames = append(funcNames, name)
	}
	sort.Strings(funcNames)
	for _, name := range funcNames {
		fn := dsl.funcs.data[name]
		// Store the parameter names for this function
		paramNames := make([]string, len(fn.meta.params))
		for i, param := range fn.meta.params {
			paramNames[i] = param.name
		}
		funcParams[name] = paramNames

		// Create a function completer that adds the complete function call with default parameters
		paramString := ""
		for i, param := range fn.meta.params {
			if i > 0 {
				paramString += " "
			}
			paramString += param.name + "="
			if param.def != nil {
				switch v := param.def.(type) {
				case string:
					// Escape any double quotes in the string and wrap in quotes
					escaped := strings.ReplaceAll(v, `"`, `\"`)
					paramString += fmt.Sprintf(`"%s"`, escaped)
				default:
					paramString += fmt.Sprintf("%v", param.def)
				}
			}
		}
		funcCompleter := readline.PcItem(name + "(" + paramString + ")")

		prefixCompleter.Children = append(prefixCompleter.Children, funcCompleter)
	}

	// Create our custom completer
	completer := &paramCompleter{
		prefixCompleter: prefixCompleter,
		funcParams:      funcParams,
	}

	// Create readline instance with basic configuration
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\x1b[34m┃ \x1b[0m",
		HistoryFile:     histFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		// Read input using readline
		input, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}
			break
		}

		// Trim whitespace and check for exit command
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		if input == "store" {
			fmt.Printf("\x1b[32mStoring current state\x1b[0m\n")
			dsl.vars.storeState()
			dsl.funcs.storeState()
			continue
		}
		if input == "restore" {
			fmt.Printf("\x1b[31mRestoring previous state\x1b[0m\n")
			dsl.vars.restoreState()
			dsl.funcs.restoreState()
			continue
		}
		if input == "?" || input == "help" {
			if input == "?" {
				fmt.Println(dsl.renderMarkdownToTerminal(buf.String()))
			} else {
				fmt.Println(dsl.docText())
			}
			continue
		}
		if input == "debug" {
			debugMode = !debugMode
			if debugMode {
				fmt.Printf("\x1b[32mDebug mode is now ON\x1b[0m\n")
			} else {
				fmt.Printf("\x1b[31mDebug mode is now OFF\x1b[0m\n")
			}
			continue
		}
		if input == "export-md" {
			filename := fmt.Sprintf("%s.md", dsl.id)
			if err := flo.File(filename).StoreString(dsl.docMarkdown()); err != nil {
				fmt.Printf("\x1b[31mError: could not generate Markdown documentation: %v\x1b[0m\n", err)
			} else {
				fmt.Printf("\x1b[32mMarkdown documentation exported to %s\x1b[0m\n", filename)
			}
			continue
		}
		if input == "export-html" {
			filename := fmt.Sprintf("%s.html", dsl.id)
			if err := flo.File(filename).StoreString(dsl.docHTML()); err != nil {
				fmt.Printf("\x1b[31mError: could not generate HTML documentation: %v\x1b[0m\n", err)
			} else {
				fmt.Printf("\x1b[32mHTML documentation exported to %s\x1b[0m\n", filename)
			}
			continue
		}
		if input == "export-vscode-extension" {
			filename, _ := filepath.Abs(fmt.Sprintf("%s.vsix", dsl.id))
			if err := dsl.exportVSCodeExtension(filename); err != nil {
				fmt.Printf("\x1b[31mError: could not generate VSCode extension: %v\x1b[0m\n", err)
			} else {
				fmt.Printf("\x1b[32mVSCode extension exported to %s\x1b[0m\n", filename)
			}
			continue
		}
		if strings.HasPrefix(input, "search ") || input == "search" {
			query := strings.TrimSpace(strings.TrimPrefix(input, "search"))
			found := false

			// Create template data
			type SearchResult struct {
				Query     string
				Found     bool
				Variables []struct {
					Name        string
					Type        string
					Description string
					Default     any
				}
				Functions []struct {
					Name        string
					Description string
					Parameters  []struct {
						Name        string
						Type        string
						Description string
						Default     any
						Min         any
						Max         any
						Unit        string
					}
					Returns []struct {
						Name        string
						Type        string
						Description string
						Default     any
						Min         any
						Max         any
						Unit        string
					}
				}
			}

			data := SearchResult{
				Query: query,
			}

			// Search variables
			for name, v := range dsl.vars.data {
				if query == "" || strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
					data.Variables = append(data.Variables, struct {
						Name        string
						Type        string
						Description string
						Default     any
					}{
						Name:        name,
						Type:        v.meta.typ,
						Description: v.meta.desc,
						Default:     v.meta.def,
					})
					found = true
				}
			}

			// Search functions
			for name, fn := range dsl.funcs.data {
				if query == "" || strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
					funcData := struct {
						Name        string
						Description string
						Parameters  []struct {
							Name        string
							Type        string
							Description string
							Default     any
							Min         any
							Max         any
							Unit        string
						}
						Returns []struct {
							Name        string
							Type        string
							Description string
							Default     any
							Min         any
							Max         any
							Unit        string
						}
					}{
						Name:        name,
						Description: fn.meta.desc,
					}

					for _, p := range fn.meta.params {
						funcData.Parameters = append(funcData.Parameters, struct {
							Name        string
							Type        string
							Description string
							Default     any
							Min         any
							Max         any
							Unit        string
						}{
							Name:        p.name,
							Type:        p.typ,
							Description: p.desc,
							Default:     p.def,
							Min:         p.min,
							Max:         p.max,
							Unit:        p.unit,
						})
					}

					for _, r := range fn.meta.returns {
						funcData.Returns = append(funcData.Returns, struct {
							Name        string
							Type        string
							Description string
							Default     any
							Min         any
							Max         any
							Unit        string
						}{
							Name:        r.name,
							Type:        r.typ,
							Description: r.desc,
							Default:     r.def,
							Min:         r.min,
							Max:         r.max,
							Unit:        r.unit,
						})
					}

					data.Functions = append(data.Functions, funcData)
					found = true
				}
			}

			data.Found = found

			// Parse and execute template
			tmpl, err := template.ParseFS(dslTemplates, "template_search.tmpl")
			if err != nil {
				fmt.Printf("\x1b[31mError parsing search template: %v\x1b[0m\n", err)
				continue
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				fmt.Printf("\x1b[31mError executing search template: %v\x1b[0m\n", err)
				continue
			}

			fmt.Println(dsl.renderMarkdownToTerminal(buf.String()))
			continue
		}

		// Execute the input
		result, err := dsl.run(input, debugMode)
		if err != nil {
			fmt.Printf("\x1b[31mError: %v\x1b[0m\n", err)
			continue
		}
		if result == nil {
			fmt.Println("No result")
			continue
		}

		if result.err != nil {
			fmt.Printf("\x1b[31m┃ Error: %v\x1b[0m\n", result.err)
			continue
		}

		// Print the result
		resStr := ""
		switch result.value.(type) {
		case color.RGBA, color.RGBA64, color.NRGBA, color.NRGBA64:
			resStr = dsl.shellResultColor(result.value.(color.Color))
		case *image.RGBA, *image.NRGBA, *image.RGBA64, *image.NRGBA64:
			resStr = dsl.shellResultImage(result.value.(image.Image))
		case Point:
			resStr = dsl.shellResultPoint(result.value.(Point))
		case Rect:
			resStr = dsl.shellResultRect(result.value.(Rect))
		// TODO: NEW TYPES: add additional types
		default:
			resStr = fmt.Sprint(result.value)
		}
		fmt.Printf("\x1b[32m┃ %v\x1b[0m\n", resStr)
	}
}
