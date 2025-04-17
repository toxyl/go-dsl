package main

import "fmt"

// dslTokenType represents the type of a token in the language.
// Tokens are the basic building blocks of the language's syntax.
type dslTokenType string

// dslNodeKind represents the type of a node in the Abstract Syntax Tree (AST).
// Each kind corresponds to a different type of language construct that can be
// parsed and evaluated.
type dslNodeKind int

func dslError(fmtStr string, args ...any) error {
	return fmt.Errorf(fmtStr, args...)
}

// dslCollection is a helper struct to group functions, variables and constants
// under a simple prefix which helps to avoid name collisions with
// existing package level functions / variables / consts.
type dslCollection struct {
	// Language metadata
	id          string // e.g. "test-script"
	name        string // e.g. "Test Script"
	description string // e.g. "A test scripting language"
	version     string // e.g. "0.0.1"
	extension   string // e.g. "ts" (without dot)
	theme       *dslColorTheme

	// Existing fields
	tokenizer *dslTokenizer
	parser    *dslParser
	ast       *dslNode
	vars      *dslVarRegistry
	funcs     *dslFnRegistry
	errors    struct {
		UNSUPPORTED_TARGET_TYPE             func(typ string) error
		STRING_CAST                         func(str, typ string) error
		NIL_CAST                            func() error
		UNSUPPORTED_SOURCE_TYPE             func(v any) error
		TKN_ASSIGN_VALUE_MISSING            func() error
		TKN_ASSIGN_NAME_MISSING             func() error
		TKN_FUNC_INCOMPLETE                 func() error
		TKN_NOT_VALID                       func(v string) error
		TKN_FUNC_WITH_SPACE                 func() error
		TKN_PAREN_MISMATCH                  func() error
		TKN_UNTERMINATED_STRING             func(pos int) error
		TKN_UNTERMINATED_COMMENT            func(pos int) error
		TKN_UNTERMINATED_FUNC               func(pos int) error
		TKN_UNTERMINATED_ARG                func(pos int) error
		TKN_ASSIGN_UNEXPECTED               func(pos int) error
		TKN_INVALID_ARG_REF                 func(pos int, reason string) error
		REG_VALIDATION_WRONG_TYPE           func(typ, name, expected string, got any) error
		REG_VALIDATION_OUT_OF_BOUNDS        func(typ, name string, min, max, got any) error
		REG_VALIDATION_OUT_OF_BOUNDS_LENGTH func(typ, name string, min, max, got any) error
		PSR_INPUT_EMPTY                     func() error
		PSR_EXPECTED_ARG                    func() error
		PSR_UNEXPECTED_TOKEN_TYPE           func(token *dslToken) error
		PSR_UNEXPECTED_OPENING_PAREN        func() error
		PSR_UNEXPECTED_CLOSING_PAREN        func() error
		PSR_ASSIGN_MISSING_NAME             func() error
		PSR_ASSIGN_MISSING_VALUE            func() error
		PSR_ASSIGN_INVALID                  func() error
		PSR_ARG_REF_INVALID                 func(ref string) error
		PSR_ARG_REF_OUT_OF_RANGE            func(id int) error
		PSR_VAR_UNDEFINED                   func(name string) error
		PSR_FUNC_UNKNOWN                    func(name string) error
		PSR_PARAM_UNKNOWN                   func(name string) error
		PSR_PARAM_STYLE_MISMATCH            func() error
		PSR_PARAM_TOO_MANY                  func(name string) error
		PSR_UNSUPPORTED_NODE_TYPE           func(node *dslNode) error
	}
	tokens struct {
		invalid    dslTokenType // Invalid token
		argRef     dslTokenType // Argument reference ($1, $2, etc.)
		argValue   dslTokenType // Argument value
		boolean    dslTokenType // Boolean literal
		comment    dslTokenType // Comment
		terminator dslTokenType // Statement terminator
		float      dslTokenType // Floating-point literal
		callEnd    dslTokenType // End of function call
		callStart  dslTokenType // Start of function call
		integer    dslTokenType // Integer literal
		namedArg   dslTokenType // Named argument
		null       dslTokenType // Nil value
		str        dslTokenType // String literal
		uinteger   dslTokenType // Unsigned integer literal
		assign     dslTokenType // Variable assignment
		varRef     dslTokenType // Variable reference
		space      dslTokenType // Whitespace
		fnName     dslTokenType // Function name
		parenOpen  dslTokenType // Open parenthesis
		parenClose dslTokenType // Close parenthesis
	}
	nodes struct {
		call       dslNodeKind // Function call node (e.g. print("hello"))
		arg        dslNodeKind // Argument node (e.g. "hello" in print("hello"))
		varRef     dslNodeKind // Variable reference node (e.g. x in print(x))
		str        dslNodeKind // String literal node (e.g. "hello world")
		float      dslNodeKind // Floating-point number node (e.g. 3.14)
		integer    dslNodeKind // Integer number node (e.g. 42)
		boolean    dslNodeKind // Boolean value node (e.g. true, false)
		assign     dslNodeKind // Variable assignment node (e.g. x: 42)
		terminator dslNodeKind // End of statement marker
		argRef     dslNodeKind // Script argument reference node (e.g. $1, $2)
	}
}

var dsl = dslCollection{
	id:          "",
	name:        "",
	description: "",
	version:     "",
	extension:   "",
	theme:       &dslColorTheme{},
	errors: struct {
		UNSUPPORTED_TARGET_TYPE             func(typ string) error
		STRING_CAST                         func(str string, typ string) error
		NIL_CAST                            func() error
		UNSUPPORTED_SOURCE_TYPE             func(v any) error
		TKN_ASSIGN_VALUE_MISSING            func() error
		TKN_ASSIGN_NAME_MISSING             func() error
		TKN_FUNC_INCOMPLETE                 func() error
		TKN_NOT_VALID                       func(v string) error
		TKN_FUNC_WITH_SPACE                 func() error
		TKN_PAREN_MISMATCH                  func() error
		TKN_UNTERMINATED_STRING             func(pos int) error
		TKN_UNTERMINATED_COMMENT            func(pos int) error
		TKN_UNTERMINATED_FUNC               func(pos int) error
		TKN_UNTERMINATED_ARG                func(pos int) error
		TKN_ASSIGN_UNEXPECTED               func(pos int) error
		TKN_INVALID_ARG_REF                 func(pos int, reason string) error
		REG_VALIDATION_WRONG_TYPE           func(typ string, name string, expected string, got any) error
		REG_VALIDATION_OUT_OF_BOUNDS        func(typ, name string, min, max, got any) error
		REG_VALIDATION_OUT_OF_BOUNDS_LENGTH func(typ, name string, min, max, got any) error
		PSR_INPUT_EMPTY                     func() error
		PSR_EXPECTED_ARG                    func() error
		PSR_UNEXPECTED_TOKEN_TYPE           func(token *dslToken) error
		PSR_UNEXPECTED_OPENING_PAREN        func() error
		PSR_UNEXPECTED_CLOSING_PAREN        func() error
		PSR_ASSIGN_MISSING_NAME             func() error
		PSR_ASSIGN_MISSING_VALUE            func() error
		PSR_ASSIGN_INVALID                  func() error
		PSR_ARG_REF_INVALID                 func(ref string) error
		PSR_ARG_REF_OUT_OF_RANGE            func(id int) error
		PSR_VAR_UNDEFINED                   func(name string) error
		PSR_FUNC_UNKNOWN                    func(name string) error
		PSR_PARAM_UNKNOWN                   func(name string) error
		PSR_PARAM_STYLE_MISMATCH            func() error
		PSR_PARAM_TOO_MANY                  func(name string) error
		PSR_UNSUPPORTED_NODE_TYPE           func(node *dslNode) error
	}{
		UNSUPPORTED_TARGET_TYPE:  func(typ string) error { return dslError("unsupported target type: %s", typ) },
		STRING_CAST:              func(str, typ string) error { return dslError("cannot cast string %q to %s", str, typ) },
		NIL_CAST:                 func() error { return dslError("cannot cast nil value") },
		UNSUPPORTED_SOURCE_TYPE:  func(v any) error { return dslError("unsupported source type: %T", v) },
		TKN_ASSIGN_VALUE_MISSING: func() error { return dslError("missing var value in assign") },
		TKN_ASSIGN_NAME_MISSING:  func() error { return dslError("missing var name in assign") },
		TKN_FUNC_INCOMPLETE:      func() error { return dslError("func call incomplete") },
		TKN_NOT_VALID:            func(v string) error { return dslError("'%s' is not a valid token", v) },
		TKN_FUNC_WITH_SPACE:      func() error { return dslError("function names cannot contain whitespaces") },
		TKN_PAREN_MISMATCH:       func() error { return dslError("parenthesis mismatch") },
		TKN_UNTERMINATED_STRING:  func(pos int) error { return dslError("unterminated string at position %d", pos) },
		TKN_UNTERMINATED_COMMENT: func(pos int) error { return dslError("unterminated comment at position %d", pos) },
		TKN_UNTERMINATED_FUNC:    func(pos int) error { return dslError("unterminated function at position %d", pos) },
		TKN_UNTERMINATED_ARG:     func(pos int) error { return dslError("unterminated argument at position %d", pos) },
		TKN_ASSIGN_UNEXPECTED:    func(pos int) error { return dslError("unexpected variable assignment at position %d", pos) },
		TKN_INVALID_ARG_REF: func(pos int, reason string) error {
			return dslError("invalid argument reference at position %d: %s", pos, reason)
		},
		REG_VALIDATION_WRONG_TYPE: func(typ, name, expected string, got any) error {
			return dslError("%s %s: expected %s, got %T", typ, name, expected, got)
		},
		REG_VALIDATION_OUT_OF_BOUNDS: func(typ, name string, min, max, got any) error {
			return dslError("%s %s: value %v is out of bounds (%v - %v)", typ, name, got, min, max)
		},
		REG_VALIDATION_OUT_OF_BOUNDS_LENGTH: func(typ, name string, min, max, got any) error {
			return dslError("%s %s: length %v is out of bounds (%v - %v)", typ, name, got, min, max)
		},
		PSR_INPUT_EMPTY:              func() error { return dslError("input is empty") },
		PSR_EXPECTED_ARG:             func() error { return dslError("expected argument") },
		PSR_UNEXPECTED_TOKEN_TYPE:    func(token *dslToken) error { return dslError("unexpected token type: %s", token.Type) },
		PSR_UNEXPECTED_OPENING_PAREN: func() error { return dslError("unexpected opening parenthesis") },
		PSR_UNEXPECTED_CLOSING_PAREN: func() error { return dslError("unexpected closing parenthesis") },
		PSR_ASSIGN_MISSING_NAME:      func() error { return dslError("missing variable name in assignment") },
		PSR_ASSIGN_MISSING_VALUE:     func() error { return dslError("expected value after variable assignment") },
		PSR_ASSIGN_INVALID:           func() error { return dslError("invalid variable assignment") },
		PSR_ARG_REF_INVALID:          func(ref string) error { return dslError("invalid argument reference: %s", ref) },
		PSR_ARG_REF_OUT_OF_RANGE:     func(id int) error { return dslError("argument $%d out of range", id) },
		PSR_VAR_UNDEFINED:            func(name string) error { return dslError("undefined variable: %s", name) },
		PSR_FUNC_UNKNOWN:             func(name string) error { return dslError("unknown function: %s", name) },
		PSR_PARAM_UNKNOWN:            func(name string) error { return dslError("unknown parameter: %s", name) },
		PSR_PARAM_STYLE_MISMATCH:     func() error { return dslError("must use positional or named arguments, not both") },
		PSR_PARAM_TOO_MANY:           func(name string) error { return dslError("too many arguments for function %s", name) },
		PSR_UNSUPPORTED_NODE_TYPE:    func(node *dslNode) error { return dslError("unsupported node type: %v", node.kind) },
	},
	tokens: struct {
		invalid    dslTokenType
		argRef     dslTokenType
		argValue   dslTokenType
		boolean    dslTokenType
		comment    dslTokenType
		terminator dslTokenType
		float      dslTokenType
		callEnd    dslTokenType
		callStart  dslTokenType
		integer    dslTokenType
		namedArg   dslTokenType
		null       dslTokenType
		str        dslTokenType
		uinteger   dslTokenType
		assign     dslTokenType
		varRef     dslTokenType
		space      dslTokenType
		fnName     dslTokenType
		parenOpen  dslTokenType
		parenClose dslTokenType
	}{
		invalid:    "INVALID",
		argRef:     "ARG_REF",
		argValue:   "VALUE",
		boolean:    "BOOL",
		comment:    "COMMENT",
		terminator: "TERMINATOR",
		float:      "FLOAT",
		callEnd:    "CALL_END",
		callStart:  "CALL_START",
		integer:    "INT",
		namedArg:   "ARG",
		null:       "NIL",
		str:        "STRING",
		uinteger:   "UINT",
		assign:     "ASSIGN",
		varRef:     "VAR",
		space:      "WHITESPACE",
		fnName:     "FUNC_NAME",
		parenOpen:  "OPEN_PAREN",
		parenClose: "CLOSE_PAREN",
	},
	nodes: struct {
		call       dslNodeKind
		arg        dslNodeKind
		varRef     dslNodeKind
		str        dslNodeKind
		float      dslNodeKind
		integer    dslNodeKind
		boolean    dslNodeKind
		assign     dslNodeKind
		terminator dslNodeKind
		argRef     dslNodeKind
	}{
		call:       0,
		arg:        1,
		varRef:     2,
		str:        3,
		float:      4,
		integer:    5,
		boolean:    6,
		assign:     7,
		terminator: 8,
		argRef:     9,
	},
}
