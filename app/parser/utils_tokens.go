package main

import (
	"slices"
	"strings"

	"github.com/toxyl/math"
)

func (dsl *dslCollection) getLastToken(tokens []*dslToken) *dslToken {
	return tokens[math.Max(0, len(tokens)-1)]
}

func (dsl *dslCollection) isAnyToken(token *dslToken, types ...dslTokenType) bool {
	return slices.Contains(types, token.Type)
}

func (dsl *dslCollection) trimTokenSpace(token *dslToken) {
	token.Value = strings.TrimSpace(token.Value)
}
func (dsl *dslCollection) containsTokenSpace(token *dslToken) bool {
	return strings.ContainsAny(token.Value, " \t\r\n")
}
func (dsl *dslCollection) trimToken(token *dslToken, cutset string) {
	token.Value = strings.Trim(token.Value, cutset)
}
func (dsl *dslCollection) trimTokenRight(token *dslToken, cutset string) {
	token.Value = strings.TrimRight(token.Value, cutset)
}

func (dsl *dslCollection) replaceInToken(token *dslToken, search, replace string) {
	token.Value = strings.ReplaceAll(token.Value, search, replace)
}

func (dsl *dslCollection) appendToken(token *dslToken, value string) { token.Value += value }

func (dsl *dslCollection) isCallEndToken(token *dslToken) bool {
	return token.Type == dsl.tokens.callEnd
}
func (dsl *dslCollection) isCallStartToken(token *dslToken) bool {
	return token.Type == dsl.tokens.callStart
}
func (dsl *dslCollection) isAssignToken(token *dslToken) bool {
	return token.Type == dsl.tokens.assign
}
func (dsl *dslCollection) isTerminatorToken(token *dslToken) bool {
	return token.Type == dsl.tokens.terminator
}
func (dsl *dslCollection) isCommentToken(token *dslToken) bool {
	return token.Type == dsl.tokens.comment
}
func (dsl *dslCollection) isStringToken(token *dslToken) bool {
	return token.Type == dsl.tokens.str
}
func (dsl *dslCollection) isArgValueToken(token *dslToken) bool {
	return token.Type == dsl.tokens.argValue
}
func (dsl *dslCollection) isInvalidToken(token *dslToken) bool {
	return token.Type == dsl.tokens.invalid
}
func (dsl *dslCollection) isEmptyToken(token *dslToken) bool { return token.Value == "" }

func (dsl *dslCollection) isNotEmptyToken(token *dslToken) bool { return token.Value != "" }
func (dsl *dslCollection) isNotStringToken(token *dslToken) bool {
	return token.Type != dsl.tokens.str
}
func (dsl *dslCollection) isNotCommentToken(token *dslToken) bool {
	return token.Type != dsl.tokens.comment
}
func (dsl *dslCollection) isNotCallStartToken(token *dslToken) bool {
	return token.Type != dsl.tokens.callStart
}
func (dsl *dslCollection) isNotAssignToken(token *dslToken) bool {
	return token.Type != dsl.tokens.assign
}
func (dsl *dslCollection) isNotTerminatorToken(token *dslToken) bool {
	return token.Type != dsl.tokens.terminator
}
