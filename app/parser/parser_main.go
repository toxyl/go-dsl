package main

import (
	"image"
	"reflect"
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
	case dsl.tokens.comment:
		return nil, nil
	case dsl.tokens.varRef:
		return &dslNode{
			kind: dsl.nodes.varRef,
			data: p.curr.Value,
		}, nil
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
	case dsl.tokens.sliceStart:
		return p.parseSlice()
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
		if p.curr.Type == dsl.tokens.comment {
			continue
		}
		if p.curr.Type == dsl.tokens.sliceEnd {
			// End of loop body reached
			break
		}

		arg, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if arg != nil {
			node.children = append(node.children, arg)
		}
	}

	return node, nil
}

// parseSlice parses a slice literal with its elements.
// It handles numeric literals, variable references, and function calls
// as slice elements. Elements are space-separated within curly braces.
// Returns an error if the slice syntax is invalid or if element parsing fails.
func (p *dslParser) parseSlice() (*dslNode, error) {
	// Two modes: flat 1D slice or angle-bracket rows -> matrix
	elements := make([]*dslNode, 0)
	rows := make([]*dslNode, 0)
	sawRow := false

	for p.advance() {
		if p.curr.Type == dsl.tokens.sliceEnd {
			break
		}
		if p.curr.Type == dsl.tokens.space || p.curr.Value == "" {
			continue
		}
		if p.curr.Type == dsl.tokens.comment {
			continue
		}
		if p.curr.Type == dsl.tokens.rowStart {
			sawRow = true
			row := &dslNode{kind: dsl.nodes.row}
			// collect row elements until rowEnd
			for p.advance() {
				if p.curr.Type == dsl.tokens.rowEnd {
					break
				}
				if p.curr.Type == dsl.tokens.space || p.curr.Value == "" {
					continue
				}
				if p.curr.Type == dsl.tokens.comment {
					continue
				}
				arg, err := p.parseNode()
				if err != nil {
					return nil, err
				}
				if arg != nil {
					row.children = append(row.children, arg)
				}
			}
			rows = append(rows, row)
			continue
		}

		// flat element
		arg, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if arg != nil {
			elements = append(elements, arg)
		}
	}

	if sawRow {
		return &dslNode{kind: dsl.nodes.matrix, children: rows}, nil
	}
	return &dslNode{kind: dsl.nodes.slice, children: elements}, nil
}

// parseForRange parses a for loop construct: for target[vars]{ body }
func (p *dslParser) parseForRange() (*dslNode, error) {
	node := &dslNode{
		kind: dsl.nodes.forRange,
	}

	if !p.advance() {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	targetName := p.curr.Value
	target := &dslNode{
		kind: dsl.nodes.varRef,
		data: targetName,
	}
	node.children = append(node.children, target)

	if !p.advance() {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	if p.curr.Type != dsl.tokens.indexStart {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	varNames := []string{}
	for p.advance() {
		if p.curr.Type == dsl.tokens.indexEnd || p.curr.Type == dsl.tokens.terminator {
			break
		}
		if p.curr.Type == dsl.tokens.space || p.curr.Value == "" {
			continue
		}
		if p.curr.Type == dsl.tokens.comment {
			continue
		}
		if p.curr.Type == dsl.tokens.varRef {
			varNames = append(varNames, p.curr.Value)
		} else {
			return nil, dsl.errors.PSR_FOR_INVALID_VARS()
		}
	}

	if len(varNames) == 0 {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	node.data = strings.Join(varNames, " ")

	bodyStatements := []*dslNode{}

	if !p.advance() {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	for {
		if p.curr == nil {
			break
		}

		if p.curr.Type == dsl.tokens.done {
			break
		}

		if p.curr.Type == dsl.tokens.space || p.curr.Value == "" {
			if !p.advance() {
				break
			}
			continue
		}
		if p.curr.Type == dsl.tokens.comment {
			if !p.advance() {
				break
			}
			continue
		}

		stmt, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			bodyStatements = append(bodyStatements, stmt)
		}

		// Advance to next token
		if !p.advance() {
			break
		}
	}

	if len(bodyStatements) == 0 {
		return nil, dsl.errors.PSR_FOR_INVALID_VARS()
	}

	node.children = append(node.children, bodyStatements...)

	return node, nil
}

// parseIndex parses one or more chained index operations on a base node.
// It expects the current token to be indexStart when called.
func (p *dslParser) parseIndex(base *dslNode) (*dslNode, error) {
	// current is indexStart '['; consume tokens until matching indexEnd, building sub-parser slice
	depth := 1
	tokens := make([]*dslToken, 0)
	for p.advance() {
		if p.curr.Type == dsl.tokens.indexStart {
			depth++
		} else if p.curr.Type == dsl.tokens.indexEnd {
			depth--
			if depth == 0 {
				break
			}
		}
		tokens = append(tokens, p.curr)
	}
	if depth != 0 {
		return nil, dsl.errors.PSR_UNEXPECTED_CLOSING_PAREN()
	}
	if len(tokens) == 0 {
		return nil, dsl.errors.PSR_EXPECTED_ARG()
	}
	// Parse one or two expressions from tokens
	sub := &dslParser{curr: nil, next: nil, prev: nil, tokens: tokens, formatted: "", types: "", pos: -1, args: []any{}}
	idxParts := make([]*dslNode, 0, 2)
	for sub.advance() {
		if sub.curr.Type == dsl.tokens.terminator || sub.curr.Value == "" || sub.curr.Type == dsl.tokens.space {
			continue
		}
		// Stop if unexpected closing encountered
		if sub.curr.Type == dsl.tokens.indexEnd {
			break
		}
		n, err := sub.parseArgument()
		if err != nil {
			return nil, err
		}
		if n != nil {
			idxParts = append(idxParts, n)
			if len(idxParts) > 2 {
				return nil, dsl.errors.PSR_PARAM_TOO_MANY("index")
			}
		}
	}
	if len(idxParts) == 0 {
		return nil, dsl.errors.PSR_EXPECTED_ARG()
	}
	children := []*dslNode{base}
	children = append(children, idxParts...)
	return &dslNode{kind: dsl.nodes.index, children: children}, nil
}

// parseNode parses a single node from the token stream.
// It handles different types of nodes based on the current token.
// Returns an error if the token sequence is invalid.
func (p *dslParser) parseNode() (*dslNode, error) {
	switch p.curr.Type {
	case dsl.tokens.comment:
		return nil, nil
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
	case dsl.tokens.forLoop:
		return p.parseForRange()
	case dsl.tokens.callStart:
		return p.parseCall()
	case dsl.tokens.sliceStart:
		return p.parseSlice()
	case dsl.tokens.sliceEnd:
		return nil, nil
	case dsl.tokens.indexEnd:
		return nil, nil
	case dsl.tokens.callEnd:
		parens := 0
		sliceDepth := 0
		for _, t := range p.tokens {
			switch t.Type {
			case dsl.tokens.sliceStart:
				sliceDepth++
			case dsl.tokens.sliceEnd:
				sliceDepth--
			case dsl.tokens.callStart:
				if sliceDepth == 0 {
					parens++
				}
			case dsl.tokens.callEnd:
				if sliceDepth == 0 {
					parens--
				}
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
			base := &dslNode{
				kind: dsl.nodes.varRef,
				data: p.curr.Value,
			}
			// if next is an indexStart, parse indexing (can be chained)
			for p.next != nil && p.next.Type == dsl.tokens.indexStart {
				// move to indexStart and parse
				if !p.advance() {
					break
				}
				if p.curr.Type != dsl.tokens.indexStart {
					break
				}
				idxNode, err := p.parseIndex(base)
				if err != nil {
					return nil, err
				}
				base = idxNode
				// after parseIndex, p.curr == indexEnd; loop will check if another [ follows
			}
			return base, nil
		}
		if p.next != nil && p.next.Type == dsl.tokens.callStart {
			return p.parseCall()
		}
		// Support indexing after a completed call expression: (handled when parseCall returns and the caller sees indexStart)
		// Allow indexing the result of a call or other expressions by handling indexStart after parseCall above
		if p.curr.Type == dsl.tokens.indexStart {
			// This path occurs when an expression directly followed by [ is being parsed inside a larger context
			// Fall back: error because we don't have a base
			return nil, dsl.errors.PSR_EXPECTED_ARG()
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
	case dsl.nodes.slice:
		// Evaluate all children first
		vals := make([]any, 0, len(node.children))
		for _, child := range node.children {
			v, err := p.evaluateNode(child)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}

		// Empty slice -> []any{}
		if len(vals) == 0 {
			return []any{}, nil
		}

		// Helper checks
		allNumeric := true
		allStrings := true
		firstNonNilType := reflect.TypeOf(vals[0])
		uniformType := true
		for _, v := range vals {
			if v == nil {
				// nil breaks numeric and string checks, and type uniformity
				allNumeric = false
				allStrings = false
				uniformType = false
				continue
			}
			// numeric check via toFloat64
			if _, err := dsl.toFloat64(v); err != nil {
				allNumeric = false
			}
			if _, ok := v.(string); !ok {
				allStrings = false
			}
			if t := reflect.TypeOf(v); firstNonNilType == nil {
				firstNonNilType = t
			} else if t != firstNonNilType {
				uniformType = false
			}
		}

		if allNumeric {
			res := make([]float64, 0, len(vals))
			for _, v := range vals {
				f, _ := dsl.toFloat64(v)
				res = append(res, f)
			}
			return res, nil
		}

		if allStrings {
			res := make([]string, 0, len(vals))
			for _, v := range vals {
				res = append(res, v.(string))
			}
			return res, nil
		}

		// If all same type, and type matches one of the supported custom types, return typed slice
		if uniformType && firstNonNilType != nil {
			switch firstNonNilType {
			case reflect.TypeOf(&Ellipse{}):
				out := make([]*Ellipse, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Ellipse))
				}
				return out, nil
			case reflect.TypeOf(&NGon{}):
				out := make([]*NGon, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*NGon))
				}
				return out, nil
			case reflect.TypeOf(&Point{}):
				out := make([]*Point, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Point))
				}
				return out, nil
			case reflect.TypeOf(&Quad{}):
				out := make([]*Quad, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Quad))
				}
				return out, nil
			case reflect.TypeOf(&Rect{}):
				out := make([]*Rect, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Rect))
				}
				return out, nil
			case reflect.TypeOf(&LineStyle{}):
				out := make([]*LineStyle, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*LineStyle))
				}
				return out, nil
			case reflect.TypeOf(&FillStyle{}):
				out := make([]*FillStyle, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*FillStyle))
				}
				return out, nil
			case reflect.TypeOf(&TextStyle{}):
				out := make([]*TextStyle, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*TextStyle))
				}
				return out, nil
			case reflect.TypeOf(&Text{}):
				out := make([]*Text, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Text))
				}
				return out, nil
			case reflect.TypeOf(&Triangle{}):
				out := make([]*Triangle, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Triangle))
				}
				return out, nil
			case reflect.TypeOf(&Vector{}):
				out := make([]*Vector, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*Vector))
				}
				return out, nil
			case reflect.TypeOf((*image.RGBA)(nil)):
				out := make([]*image.RGBA, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*image.RGBA))
				}
				return out, nil
			case reflect.TypeOf((*image.NRGBA)(nil)):
				out := make([]*image.NRGBA, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*image.NRGBA))
				}
				return out, nil
			case reflect.TypeOf((*image.RGBA64)(nil)):
				out := make([]*image.RGBA64, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*image.RGBA64))
				}
				return out, nil
			case reflect.TypeOf((*image.NRGBA64)(nil)):
				out := make([]*image.NRGBA64, 0, len(vals))
				for _, v := range vals {
					out = append(out, v.(*image.NRGBA64))
				}
				return out, nil
			}
		}

		// Fallback: []any (coerce numerics to float64 for consistency)
		res := make([]any, 0, len(vals))
		for _, v := range vals {
			if _, err := dsl.toFloat64(v); err == nil {
				f, _ := dsl.toFloat64(v)
				res = append(res, f)
				continue
			}
			res = append(res, v)
		}
		return res, nil
	case dsl.nodes.matrix:
		// Evaluate all rows and infer a common type across the matrix.
		// Rules mirror slice inference, applied to elements across rows.
		// Row lengths must match.
		var allVals [][]any
		expectedLen := -1
		for _, r := range node.children {
			if r.kind != dsl.nodes.row {
				return nil, dsl.errors.PSR_UNSUPPORTED_NODE_TYPE(r)
			}
			rowVals := make([]any, 0, len(r.children))
			for _, c := range r.children {
				v, err := p.evaluateNode(c)
				if err != nil {
					return nil, err
				}
				rowVals = append(rowVals, v)
			}
			if expectedLen == -1 {
				expectedLen = len(rowVals)
			} else if len(rowVals) != expectedLen {
				return nil, dsl.errors.REG_VALIDATION_OUT_OF_BOUNDS_LENGTH("matrix", "row", expectedLen, expectedLen, len(rowVals))
			}
			allVals = append(allVals, rowVals)
		}

		// Infer element type across entire matrix
		allNumeric := true
		allStrings := true
		var firstNonNilType reflect.Type
		uniformType := true
		for _, row := range allVals {
			for _, v := range row {
				if v == nil {
					allNumeric = false
					allStrings = false
					uniformType = false
					continue
				}
				if _, err := dsl.toFloat64(v); err != nil {
					allNumeric = false
				}
				if _, ok := v.(string); !ok {
					allStrings = false
				}
				t := reflect.TypeOf(v)
				if firstNonNilType == nil {
					firstNonNilType = t
				} else if t != firstNonNilType {
					uniformType = false
				}
			}
		}

		if allNumeric {
			out := make([][]float64, 0, len(allVals))
			for _, row := range allVals {
				rr := make([]float64, 0, len(row))
				for _, v := range row {
					f, _ := dsl.toFloat64(v)
					rr = append(rr, f)
				}
				out = append(out, rr)
			}
			return out, nil
		}
		if allStrings {
			out := make([][]string, 0, len(allVals))
			for _, row := range allVals {
				rr := make([]string, 0, len(row))
				for _, v := range row {
					rr = append(rr, v.(string))
				}
				out = append(out, rr)
			}
			return out, nil
		}
		if uniformType && firstNonNilType != nil {
			switch firstNonNilType {
			case reflect.TypeOf(&Ellipse{}):
				out := make([][]*Ellipse, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Ellipse, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Ellipse))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&NGon{}):
				out := make([][]*NGon, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*NGon, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*NGon))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Point{}):
				out := make([][]*Point, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Point, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Point))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Quad{}):
				out := make([][]*Quad, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Quad, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Quad))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Rect{}):
				out := make([][]*Rect, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Rect, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Rect))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&LineStyle{}):
				out := make([][]*LineStyle, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*LineStyle, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*LineStyle))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&FillStyle{}):
				out := make([][]*FillStyle, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*FillStyle, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*FillStyle))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&TextStyle{}):
				out := make([][]*TextStyle, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*TextStyle, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*TextStyle))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Text{}):
				out := make([][]*Text, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Text, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Text))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Triangle{}):
				out := make([][]*Triangle, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Triangle, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Triangle))
					}
					out = append(out, rr)
				}
				return out, nil
			case reflect.TypeOf(&Vector{}):
				out := make([][]*Vector, 0, len(allVals))
				for _, row := range allVals {
					rr := make([]*Vector, 0, len(row))
					for _, v := range row {
						rr = append(rr, v.(*Vector))
					}
					out = append(out, rr)
				}
				return out, nil
			}
		}

		// Fallback: [][]any with numeric elements coerced to float64
		anyRows := make([][]any, 0, len(allVals))
		for _, row := range allVals {
			rr := make([]any, 0, len(row))
			for _, v := range row {
				if _, err := dsl.toFloat64(v); err == nil {
					f, _ := dsl.toFloat64(v)
					rr = append(rr, f)
				} else {
					rr = append(rr, v)
				}
			}
			anyRows = append(anyRows, rr)
		}
		return anyRows, nil
	case dsl.nodes.row:
		// Should not evaluate rows directly outside matrix context
		return nil, dsl.errors.PSR_UNSUPPORTED_NODE_TYPE(node)
	case dsl.nodes.index:
		if len(node.children) < 2 || len(node.children) > 3 {
			return nil, dsl.errors.PSR_ASSIGN_INVALID()
		}
		baseVal, err := p.evaluateNode(node.children[0])
		if err != nil {
			return nil, err
		}
		if mat, ok := baseVal.([][]float64); ok {
			if len(node.children) != 3 {
				return nil, dsl.errors.PSR_EXPECTED_ARG()
			}
			rIdxVal, err := p.evaluateNode(node.children[1])
			if err != nil {
				return nil, err
			}
			cIdxVal, err := p.evaluateNode(node.children[2])
			if err != nil {
				return nil, err
			}
			rf, err := dsl.toFloat64(rIdxVal)
			if err != nil {
				return nil, err
			}
			cf, err := dsl.toFloat64(cIdxVal)
			if err != nil {
				return nil, err
			}
			rIdx := int(rf)
			cIdx := int(cf)
			if rIdx < 0 || rIdx >= len(mat) {
				return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(rIdx)
			}
			if cIdx < 0 || cIdx >= len(mat[rIdx]) {
				return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(cIdx)
			}
			return mat[rIdx][cIdx], nil
		}
		if slice, ok := baseVal.([]float64); ok {
			idxVal, err := p.evaluateNode(node.children[1])
			if err != nil {
				return nil, err
			}
			fidx, err := dsl.toFloat64(idxVal)
			if err != nil {
				return nil, err
			}
			i := int(fidx)
			if i < 0 || i >= len(slice) {
				return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(i)
			}
			return slice[i], nil
		}

		// Generic reflection-based indexing for any slice types (1D and 2D)
		bv := reflect.ValueOf(baseVal)
		if bv.IsValid() && bv.Kind() == reflect.Slice {
			if len(node.children) == 2 {
				idxVal, err := p.evaluateNode(node.children[1])
				if err != nil {
					return nil, err
				}
				fidx, err := dsl.toFloat64(idxVal)
				if err != nil {
					return nil, err
				}
				i := int(fidx)
				if i < 0 || i >= bv.Len() {
					return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(i)
				}
				return bv.Index(i).Interface(), nil
			}
			if len(node.children) == 3 {
				// 2D indexing: base must be slice of slices
				if bv.Len() == 0 {
					return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(0)
				}
				if bv.Type().Elem().Kind() != reflect.Slice {
					return nil, dsl.errors.CAST_NOT_POSSIBLE("index base", "[][]T")
				}
				rIdxVal, err := p.evaluateNode(node.children[1])
				if err != nil {
					return nil, err
				}
				cIdxVal, err := p.evaluateNode(node.children[2])
				if err != nil {
					return nil, err
				}
				rf, err := dsl.toFloat64(rIdxVal)
				if err != nil {
					return nil, err
				}
				cf, err := dsl.toFloat64(cIdxVal)
				if err != nil {
					return nil, err
				}
				r := int(rf)
				c := int(cf)
				if r < 0 || r >= bv.Len() {
					return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(r)
				}
				row := bv.Index(r)
				if row.Kind() != reflect.Slice {
					return nil, dsl.errors.CAST_NOT_POSSIBLE("index base", "[][]T")
				}
				if c < 0 || c >= row.Len() {
					return nil, dsl.errors.PSR_ARG_REF_OUT_OF_RANGE(c)
				}
				return row.Index(c).Interface(), nil
			}
		}
		return nil, dsl.errors.CAST_NOT_POSSIBLE("index base", "slice or slice of slices")
	case dsl.nodes.forRange:
		if len(node.children) < 2 {
			return nil, dsl.errors.PSR_FOR_INVALID_VARS()
		}

		targetVal, err := p.evaluateNode(node.children[0])
		if err != nil {
			return nil, err
		}

		varNames := strings.Fields(node.data)
		if len(varNames) == 0 {
			return nil, dsl.errors.PSR_FOR_INVALID_VARS()
		}

		target := reflect.ValueOf(targetVal)
		if !target.IsValid() || target.Kind() != reflect.Slice {
			return nil, dsl.errors.PSR_FOR_TARGET_NOT_ITERABLE()
		}

		if target.Len() == 0 {
			return nil, nil
		}

		if target.Type().Elem().Kind() == reflect.Slice {
			if len(varNames) == 3 {
				for i := 0; i < target.Len(); i++ {
					row := target.Index(i)
					for j := 0; j < row.Len(); j++ {
						item := row.Index(j).Interface()
						dsl.vars.set(varNames[0], float64(i))
						dsl.vars.set(varNames[1], float64(j))
						dsl.vars.set(varNames[2], item)

						for _, stmt := range node.children[1:] {
							_, err := p.evaluateNode(stmt)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			} else if len(varNames) == 2 {
				for i := 0; i < target.Len(); i++ {
					row := target.Index(i).Interface()
					dsl.vars.set(varNames[0], float64(i))
					dsl.vars.set(varNames[1], row)

					for _, stmt := range node.children[1:] {
						_, err := p.evaluateNode(stmt)
						if err != nil {
							return nil, err
						}
					}
				}
			} else {
				return nil, dsl.errors.PSR_FOR_INVALID_VARS()
			}
		} else {
			if len(varNames) != 2 {
				return nil, dsl.errors.PSR_FOR_INVALID_VARS()
			}
			for i := 0; i < target.Len(); i++ {
				item := target.Index(i).Interface()
				dsl.vars.set(varNames[0], float64(i))
				dsl.vars.set(varNames[1], item)

				for _, stmt := range node.children[1:] {
					_, err := p.evaluateNode(stmt)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		return nil, nil
	default:
		return nil, dsl.errors.PSR_UNSUPPORTED_NODE_TYPE(node)
	}
}
