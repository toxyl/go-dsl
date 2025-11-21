package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/toxyl/flo"
)

func (dsl *dslCollection) initDSL(id, name, description, version, extension string, theme *dslColorTheme) {
	dsl.mu.Lock()
	defer dsl.mu.Unlock()
	if theme == nil {
		theme = dsl.defaultColorTheme()
	}
	dsl.id = id
	dsl.name = name
	dsl.description = description
	dsl.version = version
	dsl.extension = extension
	dsl.theme = theme
	dsl.tokenizer = &dslTokenizer{
		source: "",
		pos:    0,
		token:  dsl.newToken("", tokens.invalid),
		state:  dsl.newState(),
		tokens: []*dslToken{},
	}
	dsl.vars = &dslVarRegistry{
		mu:   &sync.Mutex{},
		data: make(map[string]*dslMetaVarType),
		state: &dslRegistryState{
			data:      make(map[string]any),
			new:       make(map[string]any),
			mu:        &sync.Mutex{},
			protected: false,
		},
	}
	dsl.funcs = &dslFnRegistry{
		mu:   &sync.Mutex{},
		data: make(map[string]*dslFnType),
		state: &dslRegistryState{
			data:      make(map[string]any),
			new:       make(map[string]any),
			mu:        &sync.Mutex{},
			protected: false,
		},
	}
	dsl.parser = &dslParser{
		dsl:       dsl,
		curr:      nil,
		next:      nil,
		prev:      nil,
		tokens:    dsl.tokenizer.getTokens(),
		formatted: dsl.tokenizer.String(),
		types:     dsl.tokenizer.getTypes(),
		pos:       -1,
		args:      []any{},
	}
}

// load initializes the language with a new script and arguments.
// It resets the internal state of the tokenizer and parser, preparing them
// for processing the new script. The args parameter allows passing values
// that can be referenced within the script using $1, $2, etc.
func (dsl *dslCollection) load(script string, args ...any) {
	dsl.tokenizer.source = script
	dsl.tokenizer.tokens = []*dslToken{}
	dsl.tokenizer.pos = 0
	dsl.tokenizer.state = dsl.newState()
	dsl.tokenizer.token = dsl.newToken("", tokens.invalid)
	dsl.parser.pos = -1
	dsl.parser.tokens = []*dslToken{}
	dsl.parser.formatted = ""
	dsl.parser.types = ""
	dsl.parser.curr = nil
	dsl.parser.next = nil
	dsl.parser.prev = nil
	dsl.parser.args = args
}

func (dsl *dslCollection) expandIncludes(script string, baseDir string, stack map[string]struct{}) (string, error) {
	if stack == nil {
		stack = make(map[string]struct{})
	}
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(script))
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		rawLine := scanner.Text()

		includePath, ok := dsl.parseIncludeLine(rawLine)
		if !ok {
			builder.WriteString(rawLine)
			builder.WriteByte('\n')
			continue
		}

		resolvedPath, err := dsl.resolveIncludePath(includePath, baseDir)
		if err != nil {
			return "", fmt.Errorf("include %q (line %d): %w", includePath, lineNo, err)
		}
		if _, seen := stack[resolvedPath]; seen {
			return "", fmt.Errorf("include cycle detected at %s", resolvedPath)
		}

		stack[resolvedPath] = struct{}{}
		content := flo.File(resolvedPath).AsString()
		expanded, err := dsl.expandIncludes(content, filepath.Dir(resolvedPath), stack)
		delete(stack, resolvedPath)
		if err != nil {
			return "", err
		}

		builder.WriteString("# include \"")
		builder.WriteString(resolvedPath)
		builder.WriteString("\" #\n")
		builder.WriteString(expanded)
		if !strings.HasSuffix(expanded, "\n") {
			builder.WriteByte('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

// parseIncludeLine returns the quoted path from an include directive if the line
// contains one; otherwise it reports false.
func (dsl *dslCollection) parseIncludeLine(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "include ") {
		return "", false
	}
	pathLiteral := strings.TrimSpace(trimmed[len("include "):])
	if len(pathLiteral) < 2 || pathLiteral[0] != '"' || pathLiteral[len(pathLiteral)-1] != '"' {
		return "", false
	}
	return pathLiteral[1 : len(pathLiteral)-1], true
}

// resolveIncludePath normalizes an include path, resolving relatives against baseDir.
func (dsl *dslCollection) resolveIncludePath(path, baseDir string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	if baseDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		baseDir = cwd
	}
	return filepath.Clean(filepath.Join(baseDir, path)), nil
}

var (
	reMacroDef        = regexp.MustCompile(`(?ism)^\s*macro\s*([a-zA-Z0-9-]{1,})\((.*?)\)\s*\{\s*(.*?)\s*\};\s*$`) // Match: macro name(params){ body };
	reMacroInvocation = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9-]{1,})\((.*?)\)\s*\}\}`)                            // Match: {{ name(args) }}
)

// parseMacros extracts macro definitions and removes them from the script
func (dsl *dslCollection) parseMacros(script string) (string, error) {
	result := reMacroDef.ReplaceAllStringFunc(script, func(match string) string {
		// Extract groups
		submatches := reMacroDef.FindStringSubmatch(match)
		if len(submatches) != 4 {
			return match // Shouldn't happen, but be safe
		}

		macroName := submatches[1]
		paramStr := strings.TrimSpace(submatches[2])
		body := strings.TrimSpace(submatches[3])

		// Parse parameters
		var params []string
		if paramStr != "" {
			params = strings.Fields(paramStr)
		}

		// Store macro
		dsl.macros[macroName] = &dslMacro{
			name:   macroName,
			params: params,
			body:   body,
		}

		// Remove macro definition from script
		return ""
	})

	return result, nil
}

// expandMacros replaces macro invocations with their bodies
func (dsl *dslCollection) expandMacros(script string) (string, error) {
	result := script

	for {
		match := reMacroInvocation.FindStringSubmatch(result)
		if match == nil {
			break // No more matches
		}

		macroName := match[1]
		argStr := strings.TrimSpace(match[2])

		macro, exists := dsl.macros[macroName]
		if !exists {
			return "", fmt.Errorf("undefined macro: %q", macroName)
		}

		// Parse arguments (semicolon-separated)
		var args []string
		if argStr != "" {
			args = strings.Split(argStr, ";")
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}
		}

		// Validate argument count
		if len(args) != len(macro.params) {
			return "", fmt.Errorf("macro %q expects %d arguments, got %d", macroName, len(macro.params), len(args))
		}

		// Build replacement
		var replacement strings.Builder

		// For type 2 macros, prepend parameter assignments
		if len(macro.params) > 0 {
			for i, param := range macro.params {
				replacement.WriteString(param)
				replacement.WriteString(": ")
				replacement.WriteString(args[i])
				replacement.WriteString("\n")
			}
		}

		// Append the macro body unchanged
		replacement.WriteString(macro.body)

		// Replace the match
		result = strings.Replace(result, match[0], replacement.String(), 1)
	}

	return result, nil
}

// run runs a script and returns the results.
// It handles parsing, execution, and error handling.
// The debug parameter enables verbose output of the execution process.
// The args parameter allows passing arguments to the script.
func (dsl *dslCollection) run(script string, debug bool, args ...any) (*dslResult, error) {
	dsl.mu.Lock()
	defer dsl.mu.Unlock()

	dsl.macros = make(map[string]*dslMacro)

	script, err := dsl.expandIncludes(script, "", nil)
	if err != nil {
		return nil, err
	}

	script, err = dsl.parseMacros(script)
	if err != nil {
		return nil, err
	}

	script, err = dsl.expandMacros(script)
	if err != nil {
		return nil, err
	}

	dsl.trimSpace(&script)
	dsl.load(script, args...)

	if err := dsl.tokenizer.tokenize(); err != nil {
		return nil, formatErrorWithPosition(err, dsl.tokenizer.source, dsl.tokenizer.state.Line, dsl.tokenizer.state.Column)
	}

	if err := dsl.tokenizer.lex(); err != nil {
		return nil, formatErrorWithPosition(err, dsl.tokenizer.source, dsl.tokenizer.state.Line, dsl.tokenizer.state.Column)
	}

	dsl.parser.tokens = dsl.tokenizer.getTokens()
	dsl.parser.formatted = dsl.tokenizer.String()
	dsl.parser.types = dsl.tokenizer.getTypes()

	var firstNode *dslNode

	if len(dsl.parser.tokens) == 1 {
		token := dsl.parser.tokens[0]
		switch token.Type {
		case tokens.argRef:
			firstNode = &dslNode{
				kind:     nodes.argRef,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		case tokens.integer:
			firstNode = &dslNode{
				kind:     nodes.integer,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		case tokens.float:
			firstNode = &dslNode{
				kind:     nodes.float,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		case tokens.str:
			firstNode = &dslNode{
				kind:     nodes.str,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		case tokens.boolean:
			firstNode = &dslNode{
				kind:     nodes.boolean,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		default:
			firstNode = &dslNode{
				kind:     nodes.varRef,
				data:     token.Value,
				children: []*dslNode{},
				named:    false,
				argName:  "",
			}
		}
	}

	for dsl.parser.advance() {
		if dsl.parser.curr.Type == tokens.terminator {
			continue
		}
		if dsl.parser.curr.Type == tokens.comment {
			continue
		}

		node, err := dsl.parser.parseNode()
		if err != nil {
			return nil, formatErrorWithPosition(err, dsl.tokenizer.source, dsl.tokenizer.state.Line, dsl.tokenizer.state.Column)
		}
		if node != nil {
			if firstNode == nil {
				firstNode = node
			} else {
				current := firstNode
				for current.next != nil {
					current = current.next
				}
				current.next = node
			}
		}
	}

	ast := firstNode

	if ast == nil {
		if len(dsl.parser.tokens) == 0 {
			return nil, fmt.Errorf("script is empty")
		}
		return nil, fmt.Errorf("no nodes to evaluate: script may be empty or contain only comments")
	}

	var result *dslResult
	for ast != nil {
		if debug {
			fmt.Println(ast.toTree())
		}
		res, err := dsl.parser.evaluateNode(ast)
		result = &dslResult{res, err}
		if err != nil {
			// Use node position if available, otherwise fall back to tokenizer state
			line, col := dsl.tokenizer.state.Line, dsl.tokenizer.state.Column
			if ast.Line > 0 {
				line, col = ast.Line, ast.Column
			}
			result.err = formatErrorWithPosition(err, dsl.tokenizer.source, line, col)
			break
		}
		ast = ast.next
	}

	if result == nil {
		return nil, fmt.Errorf("no result from evaluation")
	}
	return result, result.err
}

func (dsl *dslCollection) storeState() {
	dsl.mu.Lock()
	defer dsl.mu.Unlock()
	dsl.vars.storeState()
	dsl.funcs.storeState()
}

func (dsl *dslCollection) restoreState() {
	dsl.mu.Lock()
	defer dsl.mu.Unlock()
	dsl.vars.restoreState()
	dsl.funcs.restoreState()
}
