package main

import (
	"fmt"
)

// dslTokenizer converts source code into tokens.
// It maintains parsing state and handles lexical analysis.
type dslTokenizer struct {
	source string             // Source code to tokenize
	pos    int                // Current position in source
	token  *dslToken          // Current token being built
	tokens []*dslToken        // All tokens found
	state  *dslTokenizerState // Current tokenization state
}

func (t *dslTokenizer) lex() error {
	if len(t.tokens) == 1 {
		token := t.tokens[0]
		switch token.Type {
		case dsl.tokens.assign:
			return dsl.errors.TKN_ASSIGN_VALUE_MISSING()
		case dsl.tokens.callStart, dsl.tokens.callEnd:
			return dsl.errors.TKN_FUNC_INCOMPLETE()
		case dsl.tokens.argRef, dsl.tokens.str, dsl.tokens.comment, dsl.tokens.integer, dsl.tokens.float, dsl.tokens.boolean, dsl.tokens.null:
			return nil
		}
		// this might just be a primitive, let's determine its type and return
		token.Type = dsl.tokens.invalid
		t.determineTokenType(token)
		if token.Type == dsl.tokens.invalid {
			return dsl.errors.TKN_NOT_VALID(token.Value)
		}
		return nil
	}

	// properly lex the result:
	parens := 0
	for i, token := range t.tokens {
		switch token.Type {
		case dsl.tokens.callStart:
			dsl.trimTokenSpace(token)
			if dsl.containsTokenSpace(token) {
				return dsl.errors.TKN_FUNC_WITH_SPACE()
			}
			parens++
		case dsl.tokens.callEnd:
			parens--
			if parens < 0 {
				return dsl.errors.TKN_PAREN_MISMATCH()
			}
		case dsl.tokens.assign:
			if dsl.isAssignToken(token) && dsl.isAssign(token.Value[0]) {
				return dsl.errors.TKN_ASSIGN_NAME_MISSING()
			}
			if !t.hasTokens() || dsl.isTerminatorToken(t.tokens[i+1]) {
				return dsl.errors.TKN_ASSIGN_VALUE_MISSING()
			}
		}
	}
	if parens != 0 {
		return dsl.errors.TKN_PAREN_MISMATCH()
	}
	return nil
}

// getTokens returns all tokens found during tokenization.
func (t *dslTokenizer) getTokens() []*dslToken {
	return t.tokens
}

// getTypes returns a string representation of all token types.
// Used for debugging purposes.
func (t *dslTokenizer) getTypes() string {
	var tokens []string
	for _, token := range t.tokens {
		tokens = append(tokens, fmt.Sprintf("%s{`%s`}", string(token.Type), token.Value))
	}
	return dsl.joinSpace(tokens)
}

func (t *dslTokenizer) getPrevToken(i int) *dslToken {
	if i > 0 {
		return t.tokens[i-1]
	}
	return nil
}

// String returns a formatted string representation of all tokens.
// Used for debugging purposes.
func (t *dslTokenizer) String() string {
	var tokens []string

	for i, token := range t.tokens {
		str := token.String()
		if str == ";" {
			str = ";\n"
		} else if dsl.lastCharIs(str, ':') {
			str += " "
		} else if dsl.lastCharIs(str, '=') {
			str = " " + str
		}

		prev := t.getPrevToken(i - 1)
		if dsl.isStringToken(token) {
			str = dsl.wrapString(str) + ` `
		} else if dsl.isCommentToken(token) {
			str = dsl.wrapComment(str) + ` `
		} else if dsl.isAnyToken(token, dsl.tokens.argValue, dsl.tokens.float, dsl.tokens.integer, dsl.tokens.boolean, dsl.tokens.null) {
			str = str + ` `
		} else if dsl.isCallStartToken(token) && prev != nil && dsl.isCallEndToken(prev) {
			dsl.setLastString(&tokens, ") ") // adds padding when two or more function calls are used in sequence as arguments (e.g. `add(sub(5 3) sub(3 5))`)
		} else if dsl.isCallEndToken(token) && prev != nil && dsl.isNotCallStartToken(prev) {
			dsl.trimLastStringRight(&tokens, " ") // removes padding after last argument
		} else if dsl.isCallStartToken(token) && prev != nil && dsl.isAnyToken(prev, dsl.tokens.varRef, dsl.tokens.integer, dsl.tokens.float, dsl.tokens.boolean, dsl.tokens.str, dsl.tokens.comment, dsl.tokens.null, dsl.tokens.argValue) {
			dsl.appendLastString(&tokens, " ") // adds before function call
		} else if dsl.isTerminatorToken(token) && len(tokens) > 0 {
			dsl.trimLastStringRight(&tokens, " ")
		}
		tokens = append(tokens, str)
	}
	return dsl.join(tokens)
}

func (t *dslTokenizer) hasTokens() bool {
	return len(t.tokens) > 0
}

func (t *dslTokenizer) hasCharacterLeft() bool {
	return t.pos < len(t.source)
}

func (t *dslTokenizer) hasNext() bool {
	return t.pos < len(t.source)-1
}

// addToken adds a new token to the token stream.
func (t *dslTokenizer) addToken(token dslToken) {
	if dsl.isEmpty(token.Value) || dsl.isNewline(token.Value) {
		return
	}
	if dsl.isNotStringToken(&token) && dsl.isNotCommentToken(&token) {
		dsl.replaceInToken(&token, "\n", " ")
		dsl.replaceInToken(&token, "\r", " ")
		dsl.replaceInToken(&token, "\t", " ")
		dsl.trimTokenSpace(&token)
	}

	t.tokens = append(t.tokens, &token)
}

// addTokenAndSetNext adds a token and prepares for the next token.
func (t *dslTokenizer) addTokenAndSetNext(token *dslToken, typ dslTokenType) {
	if t.hasTokens() && dsl.isTerminatorToken(token) && dsl.isTerminatorToken(dsl.getLastToken(t.tokens)) {
		return
	}
	if t.hasTokens() && dsl.isCallStartToken(token) && t.state.notInInParens() && dsl.isNotTerminatorToken(dsl.getLastToken(t.tokens)) && dsl.isNotAssignToken(dsl.getLastToken(t.tokens)) {
		t.addToken(*dsl.newTerminatorToken())
	}
	t.determineTokenType(token)
	dsl.trimTokenRight(token, " ")

	t.addToken(*token)
	(*token) = *dsl.newToken("", typ)
}

// addCallEndToken handles the end of a function call.
// Returns true if processing should continue, false if it should stop.
func (t *dslTokenizer) addCallEndToken(token *dslToken) (cont bool, err error) {
	t.addTokenAndSetNext(token, dsl.tokens.callEnd)
	token.Value = ")"
	t.state.argValueEnd()
	t.state.parenClose()
	if t.state.parensMismatch() {
		return false, dsl.errors.TKN_PAREN_MISMATCH()
	}
	if t.state.notInInParens() {
		t.state.callEnd()
		t.state.statementEnd()
		t.addTokenAndSetNext(token, dsl.tokens.terminator)
		return true, nil
	}
	t.state.argValueStart()
	return false, nil
}

// determineTokenType sets the type of a token based on its value and context.
func (t *dslTokenizer) determineTokenType(token *dslToken) {
	if dsl.isArgValueToken(token) || dsl.isInvalidToken(token) {
		v := token.Value
		switch {
		case dsl.equals(v, "true"), dsl.equals(v, "false"):
			token.Type = dsl.tokens.boolean
		case dsl.equals(v, "nil"):
			token.Type = dsl.tokens.null
		case dsl.contains(v, "."):
			token.Type = dsl.tokens.float
		case v == "":
			token.Type = dsl.tokens.str
		default:
			// this might be an int, or it's a variable, so let's check
			if dsl.onlyDigits(v) {
				token.Type = dsl.tokens.integer
			} else {
				token.Type = dsl.tokens.varRef
			}
		}
	}
}

// handleString processes string literals, handling escape sequences.
func (t *dslTokenizer) handleString() error {
	// Start of string
	t.state.stringStart()
	t.token.Type = dsl.tokens.str

	// Skip the opening quote
	t.pos++

	for t.hasCharacterLeft() {
		c := t.source[t.pos]

		// Handle escape sequences
		if dsl.isEscape(c) {
			t.state.escapeStart()
			t.pos++
			continue
		}

		// Handle closing quote if not escaped
		if dsl.isString(c) && t.state.notInEscape() {
			t.state.stringEnd()
			t.state.escapeEnd()
			t.addTokenAndSetNext(t.token, dsl.tokens.argValue)
			t.pos++
			return nil
		}

		// Add character to token value
		t.token.append(c)
		t.state.escapeEnd()
		t.pos++
	}

	// If we reach here, we hit EOF before finding the closing quote
	return dsl.errors.TKN_UNTERMINATED_STRING(t.pos)
}

// handleNamedArg processes named arguments in function calls, handling both the argument name
// and its value, while maintaining proper state for argument processing.
func (t *dslTokenizer) handleNamedArg(token *dslToken) *dslToken {
	if dsl.isCallStartToken(token) {
		t.addTokenAndSetNext(token, dsl.tokens.argValue)
		t.pos++
		return token
	}
	needToGoBack := token.Value == ""

	if needToGoBack {
		t.token = dsl.getLastToken(t.tokens)
		token = t.token
	}

	token.Type = dsl.tokens.namedArg
	dsl.trimToken(token, " ")
	dsl.appendToken(token, "=")

	// Reset inArgValue before processing the named argument
	// This ensures clean state for the next token
	wasInArgValue := t.state.inArgValue()
	t.state.argValueEnd()

	if needToGoBack {
		t.tokens = t.tokens[:len(t.tokens)-1]
	}
	t.addTokenAndSetNext(token, dsl.tokens.argValue)

	// Restore inArgValue state after processing the named argument
	t.state.setArgValue(wasInArgValue)
	t.pos++
	return token
}

// handleTerminator processes the end of a statement, performing validation checks
// for unterminated strings, comments, functions, and arguments, and resetting the
// statement state for the next statement.
func (t *dslTokenizer) handleTerminator() error {
	t.addTokenAndSetNext(dsl.newTerminatorToken(), dsl.tokens.invalid)
	t.state.statementStart()
	t.state.assignEnd()
	if t.state.inString() {
		return dsl.errors.TKN_UNTERMINATED_STRING(t.pos)
	}
	if t.state.inComment() {
		return dsl.errors.TKN_UNTERMINATED_COMMENT(t.pos)
	}
	if t.state.inCall() {
		return dsl.errors.TKN_UNTERMINATED_FUNC(t.pos)
	}
	if t.state.inArgValue() {
		return dsl.errors.TKN_UNTERMINATED_ARG(t.pos)
	}
	t.pos++
	return nil
}

// handleComment processes comments delimited by # characters, handling both comment
// start/end markers and escape sequences within comments using backslash.
func (t *dslTokenizer) handleComment(c byte, token *dslToken) bool {
	// determine if it's a comment character and not an escape character
	// for comments, i.e. "# this is a comment"
	if dsl.isComment(c) && t.state.notInEscape() {
		t.state.commentToggle()
		if t.state.inCode() { // comment token finished
			t.token.append(c)
			t.token.Type = dsl.tokens.comment
			dsl.trimToken(t.token, "# ")
			t.addTokenAndSetNext(token, dsl.tokens.invalid)
			t.pos++
			return true
		}
	}

	// if we're in a comment, add the character to the token
	if t.state.inComment() {
		t.token.append(c)
		t.state.escapeEnd()
		t.pos++
		return true
	}

	return false
}

// handleWhitespace processes whitespace characters within strings or comments,
// preserving them as part of the token content.
func (t *dslTokenizer) handleWhitespace(c byte) bool {
	// if it's a whitespace character and we're in a string or comment, add it to the token
	if dsl.isWhitespace(c) && (t.state.inString() || t.state.inComment()) {
		t.token.append(c)
		t.state.escapeEnd()
		t.pos++
		return true
	}
	return false
}

// handleEscape processes escape sequences in strings and comments.
// Supported escape sequences include: \" for quotes, \# for comment markers,
// and \\ for backslashes.
func (t *dslTokenizer) handleEscape(c byte) bool {
	if dsl.isEscape(c) && (t.state.inString() || t.state.inComment()) {
		t.state.escapeStart()
		t.pos++
		return true
	}
	return false
}

// handleCallWithoutArgs processes function calls that have no arguments.
// It handles the function name token, validates it, and prepares the tokenizer
// state for potential argument processing in subsequent tokens.
func (t *dslTokenizer) handleCallWithoutArgs(token *dslToken) error {
	if cont, err := t.addCallEndToken(token); err != nil {
		return err
	} else if cont {
		t.pos++
		return nil
	}
	t.addTokenAndSetNext(token, dsl.tokens.invalid)
	t.pos++
	return nil
}

// isTerminator checks if the current position marks the end of a statement.
// A statement ends when the tokenizer reaches the end
// of the input source, provided we're not inside a string, comment, or function call.
func (t *dslTokenizer) isTerminator() bool {
	return !t.state.inStatement() && t.state.notInString() && t.state.inCode() && t.state.notInCall() && t.state.notInArgValue()
}

// handleArgRef processes script argument references in the format $1, $2, etc.
// It validates that the argument number is a positive integer and creates a token
// for the argument reference. Returns an error if the argument number is invalid
// or if the reference appears in an invalid context (e.g., inside a string).
func (t *dslTokenizer) handleArgRef() error {
	// Skip the '$' character
	t.pos++

	// Collect the argument number
	argNum := ""
	for t.hasCharacterLeft() {
		c := t.source[t.pos]
		if dsl.isDigit(c) {
			argNum += string(c)
			t.pos++
		} else {
			break
		}
	}

	if argNum == "" {
		return dsl.errors.TKN_INVALID_ARG_REF(t.pos, "missing number after $")
	}

	t.token.Type = dsl.tokens.argRef
	t.token.Value = "$" + argNum

	// If we're in a function call or after a variable assignment, treat this as a value
	if t.state.notInCall() {
		t.addTokenAndSetNext(t.token, dsl.tokens.argValue)
	} else if t.state.inAssign() {
		// possible but ends the statement
		if t.hasCharacterLeft() {
			t.addTokenAndSetNext(t.token, dsl.tokens.terminator)
			t.addTokenAndSetNext(dsl.newTerminatorToken(), dsl.tokens.terminator)
		}
		t.addTokenAndSetNext(t.token, dsl.tokens.invalid)
	} else {
		t.addTokenAndSetNext(t.token, dsl.tokens.invalid)
	}
	return nil
}

// tokenize performs the main tokenization process.
// It converts source code into a stream of tokens.
func (t *dslTokenizer) tokenize() error {
	var (
		token = t.token
	)

	for t.hasCharacterLeft() {
		if t.isTerminator() {
			if err := t.handleTerminator(); err != nil {
				return err
			}
			continue
		}

		c := t.source[t.pos]

		// Check for argument reference
		if dsl.isArgRef(c) && t.state.notInString() && t.state.inCode() {
			if err := t.handleArgRef(); err != nil {
				return err
			}
			continue
		}

		// determine if it's an escape character
		// for strings or comments, i.e. "hello \"world\"" or "# this is a comment\# with escape #"
		if t.handleEscape(c) {
			continue
		}

		// if it's a whitespace character and we're in a string or comment, add it to the token
		if t.handleWhitespace(c) {
			continue
		}

		if t.handleComment(c, token) {
			continue
		}

		// determine if it's a string character and not an escape character
		// for strings, i.e. "hello \"world\""
		if dsl.isString(c) && t.state.notInEscape() {
			if err := t.handleString(); err != nil {
				return err
			}
			continue
		}

		if t.state.inAssign() {
			if dsl.isWhitespace(c) {
				t.pos++
				continue // eat all whitespace following the variable assignment
			}
			t.state.assignEnd()
		}

		// check if we're in a function call and we're waiting for arguments
		if t.state.waitingForArgs() {
			// check if it's a function call without args
			if dsl.isEmptyToken(token) && dsl.isCallEnd(c) {
				if t.getPrevToken(len(t.tokens)).Type == dsl.tokens.namedArg {
					// special case where a function has a named argument that is an empty string
					t.determineTokenType(token)
					t.tokens = append(t.tokens, token)
					t.tokens = append(t.tokens, dsl.newToken(")", dsl.tokens.callEnd))
					t.pos++
					continue
				} else {
					if err := t.handleCallWithoutArgs(token); err != nil {
						return err
					}
					continue
				}
			}

			if dsl.isCallEnd(c) {
				// we finished the last arg, add it
				if cont, err := t.addCallEndToken(token); err != nil {
					return err
				} else if cont {
					t.pos++
					continue
				}
				t.addTokenAndSetNext(token, dsl.tokens.argValue)
				t.pos++
				continue
			}

			// check if it's a named argument
			if dsl.isNamedArg(c) {
				token = t.handleNamedArg(token)
				continue
			}

			// check if it's an argument separator
			if dsl.isWhitespace(c) {
				dsl.trimTokenSpace(t.token)
				t.addTokenAndSetNext(token, dsl.tokens.argValue)
				t.pos++
				continue
			}

		}

		// determine if it's a variable assignment character
		// for variable assignments, i.e. "x: 1"
		if dsl.isAssign(c) {
			t.token.append(c)
			if t.state.inAssign() {
				return dsl.errors.TKN_ASSIGN_UNEXPECTED(t.pos)
			}
			t.state.assignStart()
			token.Type = dsl.tokens.assign
			if len(t.tokens) > 0 {
				t.addTokenAndSetNext(dsl.newTerminatorToken(), dsl.tokens.assign)
			}
			t.addTokenAndSetNext(token, dsl.tokens.argValue)
			t.pos++
			for t.hasNext() && dsl.isWhitespace(t.source[t.pos]) {
				t.pos++
			}
			continue
		}

		// determine if it's a function call character
		// for function calls, i.e. "func(x)"
		if dsl.isCallStart(c) {
			t.token.append(c)
			token.Type = dsl.tokens.callStart
			t.addTokenAndSetNext(token, dsl.tokens.argValue)
			t.state.parenOpen()
			t.state.callStart()
			t.state.argValueStart()
			t.pos++
			continue
		}

		// determine if it's a function call end character
		// for function calls, i.e. "func(x)"
		if dsl.isCallEnd(c) {
			if cont, err := t.addCallEndToken(token); err != nil {
				return err
			} else if cont {
				t.pos++
				continue
			}
			t.addTokenAndSetNext(token, dsl.tokens.invalid)
			t.pos++
			continue
		}

		if dsl.isTerminator(c) {
			t.token.append(c)
			t.token.Type = dsl.tokens.terminator
			t.addTokenAndSetNext(token, dsl.tokens.invalid)
			t.pos++
			continue
		}

		if dsl.isWhitespace(c) {
			t.determineTokenType(token)
			t.addTokenAndSetNext(token, dsl.tokens.invalid)
			t.pos++
			continue
		}

		// add the character to the token
		t.token.append(c)
		t.pos++
	}

	// If we have a pending token, add it
	if dsl.isNotEmptyToken(t.token) {
		t.determineTokenType(t.token)
		t.addToken(*t.token)
	}

	return nil
}
