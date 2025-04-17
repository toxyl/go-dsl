package main

import (
	"strconv"
	"strings"
)

type dslResult struct {
	value any   // The computed value
	err   error // Any error that occurred
}

// dslParser is the main dslParser type that converts tokens into an AST.
// It maintains the current position in the token stream and handles parsing state.
type dslParser struct {
	curr      *dslToken   // Current token being processed
	next      *dslToken   // Next token to be processed
	prev      *dslToken   // Previously processed token
	tokens    []*dslToken // All tokens to be processed
	pos       int         // Current position in token stream
	formatted string      // Formatted source code
	types     string      // Token types for debugging
	args      []any       // Script arguments
}

// advance advances the parser to the next token.
// Returns false if there are no more tokens to process.
func (p *dslParser) advance() (hasMore bool) {
	p.prev = p.curr
	p.next = nil
	p.pos++

	if p.pos >= len(p.tokens) {
		return false
	}

	p.curr = p.tokens[p.pos]
	if p.pos+1 < len(p.tokens) {
		p.next = p.tokens[p.pos+1]
	}
	return true
}

// parseArgument parses an argument from the token stream.
// It handles literals, variable references, and nested function calls.
// Returns an error if the argument syntax is invalid.
func (p *dslParser) parseArgument() (*dslNode, error) {
	if p.curr == nil {
		return nil, dsl.errors.PSR_EXPECTED_ARG()
	}

	switch p.curr.Type {
	case dsl.tokens.argRef:
		return &dslNode{
			kind: dsl.nodes.argRef,
			data: p.curr.Value,
		}, nil
	case dsl.tokens.integer:
		val, err := strconv.Atoi(p.curr.Value)
		if err != nil {
			return nil, err
		}
		return &dslNode{
			kind: dsl.nodes.integer,
			data: strconv.Itoa(val),
		}, nil
	case dsl.tokens.float:
		val, err := strconv.ParseFloat(p.curr.Value, 64)
		if err != nil {
			return nil, err
		}
		return &dslNode{
			kind: dsl.nodes.float,
			data: strconv.FormatFloat(val, 'f', -1, 64),
		}, nil
	case dsl.tokens.str:
		return &dslNode{
			kind: dsl.nodes.str,
			data: p.curr.Value,
		}, nil
	case dsl.tokens.boolean:
		val, err := strconv.ParseBool(p.curr.Value)
		if err != nil {
			return nil, err
		}
		return &dslNode{
			kind: dsl.nodes.boolean,
			data: strconv.FormatBool(val),
		}, nil
	case dsl.tokens.null:
		return &dslNode{
			kind: dsl.nodes.arg,
			data: "nil",
		}, nil
	case dsl.tokens.callStart:
		return p.parseCall()
	default:
		return nil, dsl.errors.PSR_UNEXPECTED_TOKEN_TYPE(p.curr)
	}
}

// parseCall parses a function call and its arguments.
// It handles both named arguments (param=value) and positional arguments,
// supporting nested function calls and various argument types.
// Returns an error if the function call syntax is invalid or if argument
// parsing fails.
func (p *dslParser) parseCall() (*dslNode, error) {
	node := &dslNode{
		kind: dsl.nodes.call,
		data: strings.TrimSuffix(p.curr.Value, "("),
	}

	// Parse arguments
	for p.advance() {
		if p.curr.Type == dsl.tokens.callEnd {
			break
		}

		arg, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		node.children = append(node.children, arg)
	}

	return node, nil
}

// parseNode parses a single node from the token stream.
// It handles different types of nodes based on the current token.
// Returns an error if the token sequence is invalid.
func (p *dslParser) parseNode() (*dslNode, error) {
	switch p.curr.Type {
	case dsl.tokens.str:
		return &dslNode{
			kind: dsl.nodes.str,
			data: p.curr.Value,
		}, nil
	case dsl.tokens.argRef:
		return &dslNode{
			kind: dsl.nodes.argRef,
			data: p.curr.Value,
		}, nil
	case dsl.tokens.callStart:
		return p.parseCall()
	case dsl.tokens.callEnd:
		parens := 0
		for _, t := range p.tokens {
			if t.Type == dsl.tokens.callStart {
				parens++
			} else if t.Type == dsl.tokens.callEnd {
				parens--
			}
		}
		if parens < 0 {
			return nil, dsl.errors.PSR_UNEXPECTED_CLOSING_PAREN()
		} else if parens > 0 {
			return nil, dsl.errors.PSR_UNEXPECTED_OPENING_PAREN()
		}
		return nil, nil
	case dsl.tokens.terminator:
		return nil, nil
	case dsl.tokens.assign:
		// Handle variable assignment
		varName := strings.TrimRight(p.curr.Value, ":")
		if varName == "" {
			return nil, dsl.errors.PSR_ASSIGN_MISSING_NAME()
		}
		if !p.advance() {
			return nil, dsl.errors.PSR_ASSIGN_MISSING_VALUE()
		}
		value, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, dsl.errors.PSR_ASSIGN_MISSING_VALUE()
		}
		return &dslNode{
			kind:     dsl.nodes.assign,
			data:     varName,
			children: []*dslNode{value},
		}, nil
	default:
		if p.curr.Type == dsl.tokens.integer {
			if i, err := strconv.ParseInt(p.curr.Value, 10, 64); err == nil {
				return &dslNode{
					kind: dsl.nodes.integer,
					data: strconv.FormatInt(i, 10),
				}, nil
			}
		}
		if p.curr.Type == dsl.tokens.float {
			if f, err := strconv.ParseFloat(p.curr.Value, 64); err == nil {
				return &dslNode{
					kind: dsl.nodes.float,
					data: strconv.FormatFloat(f, 'f', -1, 64),
				}, nil
			}
		}
		if p.curr.Type == dsl.tokens.boolean {
			if b, err := strconv.ParseBool(p.curr.Value); err == nil {
				return &dslNode{
					kind: dsl.nodes.boolean,
					data: strconv.FormatBool(b),
				}, nil
			}
		}
		if p.curr.Type == dsl.tokens.varRef {
			return &dslNode{
				kind: dsl.nodes.varRef,
				data: p.curr.Value,
			}, nil
		}
		if p.next != nil && p.next.Type == dsl.tokens.callStart {
			return p.parseCall()
		}
		if !p.advance() {
			return nil, nil
		}

		if p.prev != nil {
			if p.prev.Type == dsl.tokens.namedArg {
				dsl.trimTokenRight(p.prev, "=")
				return &dslNode{
					kind:    dsl.nodes.arg,
					data:    p.curr.Value,
					named:   true,
					argName: p.prev.Value,
				}, nil

			}
		}

		return &dslNode{
			kind: dsl.nodes.arg,
			data: p.curr.Value,
		}, nil
	}
}

// evaluateNode evaluates a node in the AST, handling different node types:
// - Function calls: Executes the function with its arguments
// - Variable references: Retrieves the variable's value
// - Literals: Returns the literal value
// - Argument references: Retrieves script arguments ($1, $2, etc.)
// Returns an error if:
// - The function doesn't exist
// - Arguments are invalid or missing
// - Variable doesn't exist
// - Type conversion fails
// - Argument reference is invalid
func (p *dslParser) evaluateNode(node *dslNode) (any, error) {
	switch node.kind {
	case dsl.nodes.argRef:
		index, err := strconv.Atoi(strings.TrimPrefix(node.data, "$"))
		if err != nil {
			return nil, dsl.errors.PSR_ARG_REF_INVALID(node.data)
		}
		if index < 1 || index > len(p.args) {
			return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(index)
		}
		return p.args[index-1], nil
	case dsl.nodes.varRef:
		val := dsl.vars.get(node.data)
		if val == nil {
			return nil, dsl.errors.PSR_VAR_UNDEFINED(node.data)
		}
		return val.get(), nil
	case dsl.nodes.arg:
		// Create a temporary parser to parse the argument
		if node.data == "" {
			return nil, dsl.errors.PSR_INPUT_EMPTY()
		}
		tokenizer := &dslTokenizer{
			source: node.data,
			pos:    0,
			token:  &dslToken{},
			tokens: []*dslToken{},
			state:  &dslTokenizerState{},
		}
		if err := tokenizer.tokenize(); err != nil {
			return nil, err
		}

		if err := tokenizer.lex(); err != nil {
			return nil, err
		}

		parser := &dslParser{
			curr:      nil,
			next:      nil,
			prev:      nil,
			tokens:    tokenizer.getTokens(),
			formatted: tokenizer.String(),
			types:     tokenizer.getTypes(),
			pos:       -1,
			args:      []any{},
		}
		if !parser.advance() {
			return node.data, nil
		}
		argNode, err := parser.parseArgument()
		if err != nil {
			return nil, err
		}
		if argNode == nil {
			// there are no more tokens to parse
			return node.data, nil
		}
		return argNode.data, nil
	case dsl.nodes.call:
		// Evaluate all child nodes first
		args := make([]any, 0)
		fn := dsl.funcs.get(node.data)
		if fn == nil {
			return nil, dsl.errors.PSR_FUNC_UNKNOWN(node.data)
		}
		orderedArgs := make([]any, len(fn.meta.params))
		for i, param := range fn.meta.params {
			orderedArgs[i] = param.def
		}
		namedArgsMode := false
		for _, child := range node.children {
			if child.named {
				namedArgsMode = true
				// Find the parameter index by name
				found := false
				for i, param := range fn.meta.params {
					if param.name == child.argName {
						var val any
						val = child.data
						if len(child.children) > 0 {
							v, err := p.evaluateNode(child.children[0])
							if err != nil {
								return nil, err
							}
							val = v
						}
						orderedArgs[i] = val
						found = true
						break
					}
				}
				if !found {
					return nil, dsl.errors.PSR_PARAM_UNKNOWN(child.data)
				}
			} else {
				if namedArgsMode {
					return nil, dsl.errors.PSR_PARAM_STYLE_MISMATCH()
				}
				val, err := p.evaluateNode(child)
				if err != nil {
					return nil, err
				}
				args = append(args, val)
			}
		}
		// Fill in positional arguments
		for i, arg := range args {
			if i >= len(orderedArgs) {
				return nil, dsl.errors.PSR_PARAM_TOO_MANY(node.data)
			}
			orderedArgs[i] = arg
		}
		return fn.call(orderedArgs...)
	case dsl.nodes.assign:
		if len(node.children) != 1 {
			return nil, dsl.errors.PSR_ASSIGN_INVALID()
		}
		val, err := p.evaluateNode(node.children[0])
		if err != nil {
			return nil, err
		}
		dsl.vars.set(node.data, val)
		return val, nil
	case dsl.nodes.str:
		return node.data, nil
	case dsl.nodes.integer:
		return strconv.ParseInt(node.data, 10, 64)
	case dsl.nodes.float:
		return strconv.ParseFloat(node.data, 64)
	case dsl.nodes.boolean:
		return strconv.ParseBool(node.data)
	default:
		return nil, dsl.errors.PSR_UNSUPPORTED_NODE_TYPE(node)
	}
}
