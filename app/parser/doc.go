package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/toxyl/flo"
)

//go:embed *.tmpl
var dslTemplates embed.FS

// dslColorTheme represents a configurable color theme for the VSCode extension
type dslColorTheme struct {
	// Editor colors
	EditorBackground string
	EditorForeground string

	// Syntax highlighting colors
	Comments             string
	BlockComments        string
	LineComments         string
	Strings              string
	SingleQuotedStrings  string
	RuneStrings          string
	RawStrings           string
	Numbers              string
	Constants            string
	Functions            string
	VariableAssignments  string
	AssignmentOperators  string
	Types                string
	Classes              string
	Packages             string
	Variables            string
	Parameters           string
	Properties           string
	Keywords             string
	ControlKeywords      string
	Operators            string
	ArithmeticOperators  string
	ComparisonOperators  string
	AddressOperators     string
	OtherKeywords        string
	StorageTypes         string
	StorageTypeModifiers string
	StorageModifiers     string
	SupportTypes         string
	SupportFunctions     string
	SupportClasses       string
	SupportConstants     string
	SupportVariables     string
	EscapeCharacters     string
	Tags                 string
	Attributes           string
}

// defaultColorTheme returns a default color theme based on the current implementation
func (dsl *dslCollection) defaultColorTheme() *dslColorTheme {
	return &dslColorTheme{
		EditorBackground:     "#1e1e1e",
		EditorForeground:     "#d4d4d4",
		Comments:             "#6A9955",
		BlockComments:        "#808080",
		LineComments:         "#808080",
		Strings:              "#FFD4D4",
		SingleQuotedStrings:  "#FFB3B3",
		RuneStrings:          "#FF8080",
		RawStrings:           "#FF6666",
		Numbers:              "#B5CEA8",
		Constants:            "#569CD6",
		Functions:            "#FFD700",
		VariableAssignments:  "#D7BA7D",
		AssignmentOperators:  "#D4D4D4",
		Types:                "#4EC9B0",
		Classes:              "#4EC9B0",
		Packages:             "#D7BA7D",
		Variables:            "#D4E6FF",
		Parameters:           "#B5DCFE",
		Properties:           "#D4D4D4",
		Keywords:             "#569CD6",
		ControlKeywords:      "#C586C0",
		Operators:            "#9CDCFE",
		ArithmeticOperators:  "#F48771",
		ComparisonOperators:  "#DCDCAA",
		AddressOperators:     "#6A9955",
		OtherKeywords:        "#569CD6",
		StorageTypes:         "#CE9178",
		StorageTypeModifiers: "#D7BA7D",
		StorageModifiers:     "#D4D4D4",
		SupportTypes:         "#4EC9B0",
		SupportFunctions:     "#DCDCAA",
		SupportClasses:       "#4EC9B0",
		SupportConstants:     "#569CD6",
		SupportVariables:     "#9CDCFE",
		EscapeCharacters:     "#808080",
		Tags:                 "#569CD6",
		Attributes:           "#9CDCFE",
	}
}

func (dsl *dslCollection) docMarkdown() string {
	type templateData struct {
		Name      string
		Version   string
		Variables []struct {
			Name        string
			Type        string
			Default     any
			Description string
		}
		Functions []struct {
			Name        string
			Description string
			Params      []struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}
			Returns []struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}
		}
	}

	data := templateData{
		Name:    dsl.name,
		Version: dsl.version,
	}

	// Add variables
	varNames := dsl.vars.names()
	sort.Strings(varNames)
	for _, name := range varNames {
		v := dsl.vars.get(name)
		if v == nil {
			continue
		}
		data.Variables = append(data.Variables, struct {
			Name        string
			Type        string
			Default     any
			Description string
		}{
			Name:        name,
			Type:        v.meta.typ,
			Default:     v.meta.def,
			Description: v.meta.desc,
		})
	}

	// Add functions
	fnNames := dsl.funcs.names()
	sort.Strings(fnNames)
	for _, name := range fnNames {
		fn := dsl.funcs.get(name)
		if fn == nil {
			continue
		}

		funcData := struct {
			Name        string
			Description string
			Params      []struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}
			Returns []struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}
		}{
			Name:        name,
			Description: fn.meta.desc,
		}

		// Add parameters
		for _, param := range fn.meta.params {
			funcData.Params = append(funcData.Params, struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}{
				Name:        param.name,
				Type:        param.typ,
				Default:     param.def,
				Min:         param.min,
				Max:         param.max,
				Unit:        param.unit,
				Description: param.desc,
			})
		}

		// Add return values
		for _, ret := range fn.meta.returns {
			funcData.Returns = append(funcData.Returns, struct {
				Name        string
				Type        string
				Default     any
				Min         any
				Max         any
				Unit        string
				Description string
			}{
				Name:        ret.name,
				Type:        ret.typ,
				Default:     ret.def,
				Min:         ret.min,
				Max:         ret.max,
				Unit:        ret.unit,
				Description: ret.desc,
			})
		}

		data.Functions = append(data.Functions, funcData)
	}

	tmpl, err := template.ParseFS(dslTemplates, "template_markdown.tmpl")
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

func (dsl *dslCollection) docHTML() string {
	return dsl.renderMarkdownToHTML(dsl.docMarkdown())
}

func (dsl *dslCollection) docText() string {
	return dsl.renderMarkdownToTerminal(dsl.docMarkdown())
}

func (dsl *dslCollection) exportVSCodeExtension(pathToVSIXFile string) error {
	// Create fixed temporary directory for extension files
	tmpDir := filepath.Join(os.TempDir(), "vscode-extension"+dsl.id)

	// Remove the directory if it exists
	if flo.Dir(tmpDir).Exists() {
		if err := flo.Dir(tmpDir).Remove(); err != nil {
			fmt.Printf("Warning: failed to remove file %s: %v", tmpDir, err)
		}
	}

	// Create the directory
	if err := flo.Dir(tmpDir).Mkdir(0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create necessary directories
	dirs := []string{
		filepath.Join(tmpDir, "src"),
		filepath.Join(tmpDir, "syntaxes"),
		filepath.Join(tmpDir, "snippets"),
		filepath.Join(tmpDir, "out"),
		filepath.Join(tmpDir, "themes"),
	}
	for _, dir := range dirs {
		if err := flo.Dir(dir).Mkdir(0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate LICENSE file
	licenseTmpl, err := template.ParseFS(dslTemplates, "template_license.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse license template: %w", err)
	}

	var licenseBuf bytes.Buffer
	if err := licenseTmpl.Execute(&licenseBuf, nil); err != nil {
		return fmt.Errorf("failed to execute license template: %w", err)
	}

	if err := flo.File(filepath.Join(tmpDir, "LICENSE")).StoreBytes(licenseBuf.Bytes()); err != nil {
		return fmt.Errorf("failed to write LICENSE file: %w", err)
	}

	// Generate package.json
	packageJSON := map[string]any{
		"name":        dsl.id,
		"displayName": dsl.name,
		"description": dsl.description,
		"version":     dsl.version,
		"engines":     map[string]string{"vscode": "^1.96.0"},
		"categories":  []string{"Programming Languages"},
		"main":        "./out/extension.js",
		"activationEvents": []string{
			"onLanguage:" + dsl.id,
		},
		"repository": map[string]string{
			"type": "git",
			"url":  "https://github.com/toxyl/godsl",
		},
		"contributes": map[string]any{
			"languages": []map[string]any{
				{
					"id":            dsl.id,
					"aliases":       []string{dsl.id, dsl.name},
					"extensions":    []string{"." + dsl.extension},
					"configuration": "./language-configuration.json",
				},
			},
			"grammars": []map[string]any{
				{
					"language":  dsl.id,
					"scopeName": "source." + dsl.id,
					"path":      "./syntaxes/" + dsl.id + ".tmLanguage.json",
				},
			},
			"snippets": []map[string]any{
				{
					"language": dsl.id,
					"path":     "./snippets/snippets.json",
				},
			},
			"themes": []map[string]any{
				{
					"label":   dsl.name + " Theme",
					"uiTheme": "vs-dark",
					"path":    "./themes/" + dsl.id + "-color-theme.json",
				},
			},
		},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "package.json"), packageJSON); err != nil {
		return err
	}

	// Generate language configuration
	langConfig := map[string]any{
		"comments": map[string]string{
			"blockComment": "#",
		},
		"brackets": [][]string{
			{"(", ")"},
		},
		"autoClosingPairs": []map[string]string{
			{"open": "(", "close": ")"},
			{"open": "\"", "close": "\""},
		},
		"surroundingPairs": []map[string]string{
			{"open": "(", "close": ")"},
			{"open": "\"", "close": "\""},
		},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "language-configuration.json"), langConfig); err != nil {
		return err
	}

	// Generate snippets
	snippets := map[string]any{
		"Function Call": map[string]any{
			"prefix":      "func",
			"body":        []string{"${1:functionName}(${2:arg1} ${3:arg2})"},
			"description": "Create a function call",
		},
		"Variable Assignment": map[string]any{
			"prefix":      "var",
			"body":        []string{"${1:variableName}: ${2:value}"},
			"description": "Create a variable assignment",
		},
		"String Literal": map[string]any{
			"prefix":      "str",
			"body":        []string{"\"${1:text}\""},
			"description": "Create a string literal",
		},
		"Comment": map[string]any{
			"prefix":      "//",
			"body":        []string{"# ${1:comment} #"},
			"description": "Create a comment",
		},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "snippets", "snippets.json"), snippets); err != nil {
		return err
	}

	// Generate TextMate grammar
	grammar := map[string]any{
		"$schema":   "https://raw.githubusercontent.com/martinring/tmlanguage/master/tmlanguage.json",
		"name":      dsl.name,
		"scopeName": "source." + dsl.id,
		"patterns": []map[string]any{
			{
				"name":  "comment.block",
				"begin": "#",
				"end":   "#",
				"patterns": []map[string]any{
					{
						"match": "\\\\#",
						"name":  "constant.character.escape",
					},
				},
			},
			{
				"name":  "string.quoted.double",
				"begin": "\"",
				"end":   "\"",
				"patterns": []map[string]any{
					{
						"name":  "constant.character.escape",
						"match": "\\\\.",
					},
				},
			},
			{
				"name":  "constant.numeric",
				"match": "\\b\\d+(\\.\\d+)?\\b",
			},
			{
				"name":  "constant.language.boolean",
				"match": "\\b(true|false)\\b",
			},
			{
				"name":  "constant.language.null",
				"match": "\\bnil\\b",
			},
			{
				"name":  "variable.parameter",
				"match": "\\$\\d+",
			},
			{
				"name":  "variable.assign",
				"match": "\\b[a-zA-Z_][a-zA-Z0-9_]*(?=\\s*[:=])",
			},
			{
				"name":  "keyword.operator.assignment",
				"match": "[:=]",
			},
			{
				"name":  "entity.name.function",
				"match": "\\b[a-zA-Z_][a-zA-Z0-9_]*\\b(?=\\s*\\()",
			},
			{
				"name":  "variable.other",
				"match": "\\b[a-zA-Z_][a-zA-Z0-9_]*\\b",
			},
		},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "syntaxes", dsl.id+".tmLanguage.json"), grammar); err != nil {
		return fmt.Errorf("failed to write grammar file: %w", err)
	}

	// Generate color theme
	themeConfig := map[string]any{
		"$schema": "vscode://schemas/color-theme",
		"name":    dsl.name + " Theme",
		"type":    "dark",
		"colors": map[string]string{
			"editor.background": dsl.theme.EditorBackground,
			"editor.foreground": dsl.theme.EditorForeground,
		},
		"tokenColors": []map[string]any{
			{
				"name":  "Comments",
				"scope": "comment",
				"settings": map[string]string{
					"foreground": dsl.theme.Comments,
				},
			},
			{
				"name":  "Block Comments",
				"scope": "comment.block",
				"settings": map[string]string{
					"foreground": dsl.theme.BlockComments,
				},
			},
			{
				"name":  "Line Comments",
				"scope": "comment.line",
				"settings": map[string]string{
					"foreground": dsl.theme.LineComments,
				},
			},
			{
				"name":  "Strings",
				"scope": "string.quoted.double",
				"settings": map[string]string{
					"foreground": dsl.theme.Strings,
				},
			},
			{
				"name":  "Single Quoted Strings",
				"scope": "string.quoted.single",
				"settings": map[string]string{
					"foreground": dsl.theme.SingleQuotedStrings,
				},
			},
			{
				"name":  "Rune Strings",
				"scope": "string.quoted.rune",
				"settings": map[string]string{
					"foreground": dsl.theme.RuneStrings,
				},
			},
			{
				"name":  "Raw Strings",
				"scope": "string.quoted.raw",
				"settings": map[string]string{
					"foreground": dsl.theme.RawStrings,
				},
			},
			{
				"name":  "Numbers",
				"scope": "constant.numeric",
				"settings": map[string]string{
					"foreground": dsl.theme.Numbers,
				},
			},
			{
				"name":  "Constants",
				"scope": "constant.language",
				"settings": map[string]string{
					"foreground": dsl.theme.Constants,
				},
			},
			{
				"name":  "Functions",
				"scope": "entity.name.function",
				"settings": map[string]string{
					"foreground": dsl.theme.Functions,
				},
			},
			{
				"name":  "Variable Assignments",
				"scope": "variable.assign",
				"settings": map[string]string{
					"foreground": dsl.theme.VariableAssignments,
				},
			},
			{
				"name":  "Assignment Operators",
				"scope": "keyword.operator.assignment",
				"settings": map[string]string{
					"foreground": dsl.theme.AssignmentOperators,
				},
			},
			{
				"name":  "Types",
				"scope": "entity.name.type",
				"settings": map[string]string{
					"foreground": dsl.theme.Types,
				},
			},
			{
				"name":  "Classes",
				"scope": "entity.name.class",
				"settings": map[string]string{
					"foreground": dsl.theme.Classes,
				},
			},
			{
				"name":  "Packages",
				"scope": "entity.name.namespace",
				"settings": map[string]string{
					"foreground": dsl.theme.Packages,
				},
			},
			{
				"name":  "Variables",
				"scope": "variable.other",
				"settings": map[string]string{
					"foreground": dsl.theme.Variables,
				},
			},
			{
				"name":  "Parameters",
				"scope": "variable.parameter",
				"settings": map[string]string{
					"foreground": dsl.theme.Parameters,
				},
			},
			{
				"name":  "Properties",
				"scope": "variable.other.property",
				"settings": map[string]string{
					"foreground": dsl.theme.Properties,
				},
			},
			{
				"name":  "Keywords",
				"scope": "keyword",
				"settings": map[string]string{
					"foreground": dsl.theme.Keywords,
				},
			},
			{
				"name":  "Control Keywords",
				"scope": "keyword.control",
				"settings": map[string]string{
					"foreground": dsl.theme.ControlKeywords,
				},
			},
			{
				"name":  "Operators",
				"scope": "keyword.operator",
				"settings": map[string]string{
					"foreground": dsl.theme.Operators,
				},
			},
			{
				"name":  "Arithmetic Operators",
				"scope": "keyword.operator.arithmetic",
				"settings": map[string]string{
					"foreground": dsl.theme.ArithmeticOperators,
				},
			},
			{
				"name":  "Comparison Operators",
				"scope": "keyword.operator.comparison",
				"settings": map[string]string{
					"foreground": dsl.theme.ComparisonOperators,
				},
			},
			{
				"name":  "Address Operators",
				"scope": "keyword.operator.address",
				"settings": map[string]string{
					"foreground": dsl.theme.AddressOperators,
				},
			},
			{
				"name":  "Other Keywords",
				"scope": "keyword.other",
				"settings": map[string]string{
					"foreground": dsl.theme.OtherKeywords,
				},
			},
			{
				"name":  "Storage Types",
				"scope": "storage.type",
				"settings": map[string]string{
					"foreground": dsl.theme.StorageTypes,
				},
			},
			{
				"name":  "Storage Type Modifiers",
				"scope": "storage.type.modifier",
				"settings": map[string]string{
					"foreground": dsl.theme.StorageTypeModifiers,
				},
			},
			{
				"name":  "Storage Modifiers",
				"scope": "storage.modifier",
				"settings": map[string]string{
					"foreground": dsl.theme.StorageModifiers,
				},
			},
			{
				"name":  "Support Types",
				"scope": "support.type",
				"settings": map[string]string{
					"foreground": dsl.theme.SupportTypes,
				},
			},
			{
				"name":  "Support Functions",
				"scope": "support.function",
				"settings": map[string]string{
					"foreground": dsl.theme.SupportFunctions,
				},
			},
			{
				"name":  "Support Classes",
				"scope": "support.class",
				"settings": map[string]string{
					"foreground": dsl.theme.SupportClasses,
				},
			},
			{
				"name":  "Support Constants",
				"scope": "support.constant",
				"settings": map[string]string{
					"foreground": dsl.theme.SupportConstants,
				},
			},
			{
				"name":  "Support Variables",
				"scope": "support.variable",
				"settings": map[string]string{
					"foreground": dsl.theme.SupportVariables,
				},
			},
			{
				"name":  "Escape Characters",
				"scope": "constant.character.escape",
				"settings": map[string]string{
					"foreground": dsl.theme.EscapeCharacters,
				},
			},
			{
				"name":  "Tags",
				"scope": "entity.name.tag",
				"settings": map[string]string{
					"foreground": dsl.theme.Tags,
				},
			},
			{
				"name":  "Attributes",
				"scope": "entity.other.attribute-name",
				"settings": map[string]string{
					"foreground": dsl.theme.Attributes,
				},
			},
		},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "themes", dsl.id+"-color-theme.json"), themeConfig); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	// Generate completion provider
	completionProvider := fmt.Sprintf(`import * as vscode from 'vscode';

export class CustomCompletionProvider implements vscode.CompletionItemProvider {
    private functions: Map<string, vscode.CompletionItem> = new Map();
    private variables: Map<string, vscode.CompletionItem> = new Map();

    constructor() {
        this.initializeBuiltIns();
    }

    private initializeBuiltIns() {
        // Add functions
        %s

        // Add variables
        %s
    }

    private addFunction(name: string, description: string, params: { name: string, type: string, description: string }[], returnType: string) {
        const item = new vscode.CompletionItem(name, vscode.CompletionItemKind.Function);
        item.documentation = new vscode.MarkdownString(description);
        
        const paramString = params.map(p => p.name + ": " + p.type).join(" ");
        item.detail = "(" + paramString + ") -> " + returnType;
        
        const paramDocs = params.map(p => "@param " + p.name + " " + p.description).join("\n");
        item.documentation.appendCodeblock(paramDocs, "typescript");
        
        this.functions.set(name, item);
    }

    private addVariable(name: string, type: string, description: string) {
        const item = new vscode.CompletionItem(name, vscode.CompletionItemKind.Variable);
        item.documentation = new vscode.MarkdownString(description);
        item.detail = ": " + type;
        this.variables.set(name, item);
    }

    public provideCompletionItems(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context: vscode.CompletionContext
    ): vscode.ProviderResult<vscode.CompletionItem[]> {
        const line = document.lineAt(position.line).text;
        const lineUntilPosition = line.substring(0, position.character);

        // After opening parenthesis or space within parentheses, show functions
        if (lineUntilPosition.endsWith("(") || 
            (lineUntilPosition.includes("(") && lineUntilPosition.endsWith(" "))) {
            return Array.from(this.functions.values());
        }

        // After equals sign for named arguments, show variables
        if (lineUntilPosition.endsWith("=")) {
            return Array.from(this.variables.values());
        }

        // After colon for variable assignment, show all completions
        if (lineUntilPosition.endsWith(":")) {
            return [...Array.from(this.functions.values()), ...Array.from(this.variables.values())];
        }

        // Default completions
        return [...Array.from(this.functions.values()), ...Array.from(this.variables.values())];
    }

    public resolveCompletionItem(
        item: vscode.CompletionItem,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.CompletionItem> {
        return item;
    }
}`, dsl.generateFunctionCode(), dsl.generateVariableCode())

	if err := flo.File(filepath.Join(tmpDir, "src", "completionProvider.ts")).StoreString(completionProvider); err != nil {
		return fmt.Errorf("failed to write completion provider: %w", err)
	}

	// Generate extension.ts
	extensionTS := fmt.Sprintf(`import * as vscode from 'vscode';
import { CustomCompletionProvider } from './completionProvider';

export function activate(context: vscode.ExtensionContext) {
    const completionProvider = new CustomCompletionProvider();
    const completionProviderDisposable = vscode.languages.registerCompletionItemProvider(
        '%s',
        completionProvider,
        '(', ':', ' ', '='
    );

    context.subscriptions.push(completionProviderDisposable);
}

export function deactivate() {}`, dsl.id)

	if err := flo.File(filepath.Join(tmpDir, "src", "extension.ts")).StoreString(extensionTS); err != nil {
		return fmt.Errorf("failed to write extension.ts: %w", err)
	}

	// Generate tsconfig.json
	tsconfig := map[string]any{
		"compilerOptions": map[string]any{
			"module":    "commonjs",
			"target":    "ES2020",
			"outDir":    "out",
			"lib":       []string{"ES2020"},
			"sourceMap": true,
			"rootDir":   "src",
			"strict":    true,
		},
		"exclude": []string{"node_modules", ".vscode-test"},
	}

	if err := dsl.writeJSON(filepath.Join(tmpDir, "tsconfig.json"), tsconfig); err != nil {
		return err
	}

	// Install dependencies
	cmd := exec.Command("npm", "init", "-y")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize npm: %w\nOutput: %s", err, output)
	}

	cmd = exec.Command("npm", "install", "--save-dev", "typescript", "@types/vscode@1.96.0")
	cmd.Dir = tmpDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w\nOutput: %s", err, output)
	}

	// Compile TypeScript files
	cmd = exec.Command("npx", "tsc", "-p", "tsconfig.json")
	cmd.Dir = tmpDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile TypeScript: %w\nOutput: %s", err, output)
	}

	// Package the extension
	cmd = exec.Command("npx", "vsce", "package", "-o", pathToVSIXFile)
	cmd.Dir = tmpDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to package extension: %w\nOutput: %s", err, output)
	}

	return nil
}

func (dsl *dslCollection) generateFunctionCode() string {
	var sb strings.Builder
	for _, name := range dsl.funcs.names() {
		fn := dsl.funcs.get(name)
		if fn == nil {
			continue
		}

		params := make([]string, len(fn.meta.params))
		for i, param := range fn.meta.params {
			params[i] = "{ name: \"" + param.name + "\", type: \"" + param.typ + "\", description: \"" + param.desc + "\" }"
		}

		returnType := "any"
		if len(fn.meta.returns) > 0 {
			returnType = fn.meta.returns[0].typ
		}

		sb.WriteString("this.addFunction(\"" + name + "\", \"" + fn.meta.desc + "\", [" + strings.Join(params, ", ") + "], \"" + returnType + "\");\n")
	}
	return sb.String()
}

func (dsl *dslCollection) generateVariableCode() string {
	var sb strings.Builder
	for name, variable := range dsl.vars.data {
		sb.WriteString("this.addVariable(\"" + name + "\", \"" + variable.meta.typ + "\", \"" + variable.meta.desc + "\");\n")
	}
	return sb.String()
}

func (dsl *dslCollection) writeJSON(path string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := flo.File(path).StoreBytes(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	return nil
}
