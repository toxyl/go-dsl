package main

import (
	"fmt"
	"sync"
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

// run runs a script and returns the results.
// It handles parsing, execution, and error handling.
// The debug parameter enables verbose output of the execution process.
// The args parameter allows passing arguments to the script.
func (dsl *dslCollection) run(script string, debug bool, args ...any) (*dslResult, error) {
	dsl.mu.Lock()
	defer dsl.mu.Unlock()
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
