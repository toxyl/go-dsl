package main

// dslTokenizerState represents the current dslTokenizerState of the tokenizer.
// It tracks various parsing contexts like strings, comments, and function calls.
type dslTokenizerState struct {
	isInString    bool // Whether we're inside a string literal
	isInEscape    bool // Whether we're processing an escape sequence
	isInComment   bool // Whether we're inside a comment
	isInCall      bool // Whether we're inside a function call
	isInArgValue  bool // Whether we're processing an argument value
	isInStatement bool // Whether we're processing a statement
	isInAssign    bool // Whether we're processing a variable assignment
	isInForLoop   bool // Whether we're inside a for loop body
	parens        int  // Nesting level of parentheses
	slices        int  // Nesting level of slices
	indexes       int  // Nesting level of indexes
}

func (s *dslTokenizerState) parenOpen()           { s.parens++ }
func (s *dslTokenizerState) parenClose()          { s.parens-- }
func (s *dslTokenizerState) notInInParens() bool  { return s.parens == 0 }
func (s *dslTokenizerState) parensMismatch() bool { return s.parens < 0 }
func (s *dslTokenizerState) waitingForArgs() bool { return s.isInCall && s.isInArgValue }
func (s *dslTokenizerState) inCode() bool         { return !s.isInComment }
func (s *dslTokenizerState) notInString() bool    { return !s.isInString }
func (s *dslTokenizerState) notInEscape() bool    { return !s.isInEscape }
func (s *dslTokenizerState) inEscape() bool       { return s.isInEscape }
func (s *dslTokenizerState) notInCall() bool      { return !s.isInCall }
func (s *dslTokenizerState) notInArgValue() bool  { return !s.isInArgValue }
func (s *dslTokenizerState) stringStart()         { s.isInString = true }
func (s *dslTokenizerState) stringEnd()           { s.isInString = false }
func (s *dslTokenizerState) inString() bool       { return s.isInString }
func (s *dslTokenizerState) commentToggle()       { s.isInComment = !s.isInComment }
func (s *dslTokenizerState) inComment() bool      { return s.isInComment }
func (s *dslTokenizerState) statementStart()      { s.isInStatement = true }
func (s *dslTokenizerState) statementEnd()        { s.isInStatement = false }
func (s *dslTokenizerState) inStatement() bool    { return s.isInStatement }
func (s *dslTokenizerState) callStart()           { s.isInCall = true }
func (s *dslTokenizerState) callEnd()             { s.isInCall = false }
func (s *dslTokenizerState) inCall() bool         { return s.isInCall }
func (s *dslTokenizerState) escapeStart()         { s.isInEscape = true }
func (s *dslTokenizerState) escapeEnd()           { s.isInEscape = false }
func (s *dslTokenizerState) assignStart()         { s.isInAssign = true }
func (s *dslTokenizerState) assignEnd()           { s.isInAssign = false }
func (s *dslTokenizerState) inAssign() bool       { return s.isInAssign }
func (s *dslTokenizerState) argValueStart()       { s.isInArgValue = true }
func (s *dslTokenizerState) argValueEnd()         { s.isInArgValue = false }
func (s *dslTokenizerState) inArgValue() bool     { return s.isInArgValue }
func (s *dslTokenizerState) setArgValue(b bool)   { s.isInArgValue = b }
func (s *dslTokenizerState) sliceOpen()           { s.slices++ }
func (s *dslTokenizerState) sliceClose()          { s.slices-- }
func (s *dslTokenizerState) inSlice() bool        { return s.slices > 0 }
func (s *dslTokenizerState) notInSlice() bool     { return s.slices == 0 }
func (s *dslTokenizerState) indexOpen()           { s.indexes++ }
func (s *dslTokenizerState) indexClose()          { s.indexes-- }
func (s *dslTokenizerState) inIndex() bool        { return s.indexes > 0 }
func (s *dslTokenizerState) notInIndex() bool     { return s.indexes == 0 }
func (s *dslTokenizerState) forLoopStart()        { s.isInForLoop = true }
func (s *dslTokenizerState) forLoopEnd()          { s.isInForLoop = false }
func (s *dslTokenizerState) inForLoop() bool      { return s.isInForLoop }
func (s *dslTokenizerState) notInForLoop() bool   { return !s.isInForLoop }

func (dsl *dslCollection) newState() *dslTokenizerState {
	return &dslTokenizerState{
		isInString:    false,
		isInEscape:    false,
		isInComment:   false,
		isInCall:      false,
		isInArgValue:  false,
		isInStatement: true,
		isInAssign:    false,
		isInForLoop:   false,
		parens:        0,
		slices:        0,
		indexes:       0,
	}
}
