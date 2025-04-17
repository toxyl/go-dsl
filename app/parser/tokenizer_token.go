package main

// dslToken represents a lexical dslToken in the language.
// It contains the dslToken's value and type.
type dslToken struct {
	Value string       // The token's value
	Type  dslTokenType // The token's type
}

// String returns the token's value as a string.
func (t *dslToken) String() string {
	return t.Value
}

// append adds a character to the token's value.
func (t *dslToken) append(char byte) {
	t.Value += string(char)
}

// newToken creates a new token with the given value and type.
func (dsl *dslCollection) newToken(value string, tokenType dslTokenType) *dslToken {
	return &dslToken{
		Value: value,
		Type:  tokenType,
	}
}

func (dsl *dslCollection) newTerminatorToken() *dslToken {
	return dsl.newToken(";", dsl.tokens.terminator)
}
