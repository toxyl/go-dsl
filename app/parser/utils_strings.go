package main

import (
	"strings"

	"github.com/toxyl/math"
)

func (dsl *dslCollection) escapeChar(str string, c byte) string {
	return strings.ReplaceAll(str, string(c), `\`+string(c))
}

func (dsl *dslCollection) wrapString(str string) string {
	return `"` + dsl.escapeChar(str, '"') + `"`
}

func (dsl *dslCollection) wrapComment(str string) string {
	return `# ` + dsl.escapeChar(str, '#') + ` #`
}

func (dsl *dslCollection) getLastChar(str string) byte {
	if len(str) == 0 {
		return 0x00
	}
	return str[math.Max(0, len(str)-1)]
}

func (dsl *dslCollection) equals(a, b string) bool {
	return strings.EqualFold(a, b)
}

func (dsl *dslCollection) contains(str, search string) bool {
	return strings.Contains(str, search)
}

func (dsl *dslCollection) onlyDigits(str string) bool {
	for _, c := range str {
		if !dsl.isDigit(byte(c)) {
			return false
		}
	}
	return true
}

func (dsl *dslCollection) lastCharIs(str string, c byte) bool {
	return dsl.getLastChar(str) == c
}

func (dsl *dslCollection) setLastString(strs *[]string, value string) {
	i := math.Max(0, len((*strs))-1)
	(*strs)[i] = value
}

func (dsl *dslCollection) appendLastString(strs *[]string, value string) {
	i := math.Max(0, len((*strs))-1)
	(*strs)[i] = (*strs)[i] + value
}

func (dsl *dslCollection) joinSpace(str []string) string {
	return strings.Join(str, " ")
}

func (dsl *dslCollection) join(str []string) string {
	return strings.Join(str, "")
}

func (dsl *dslCollection) trimSpace(str *string) {
	(*str) = strings.TrimSpace(*str)
}

func (dsl *dslCollection) trimLastStringRight(strs *[]string, cutset string) {
	i := math.Max(0, len((*strs))-1)
	(*strs)[i] = strings.TrimRight((*strs)[i], cutset)
}

func (dsl *dslCollection) isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
func (dsl *dslCollection) isDigit(c byte) bool      { return c >= '0' && c <= '9' }
func (dsl *dslCollection) isCallStart(c byte) bool  { return c == '(' }
func (dsl *dslCollection) isCallEnd(c byte) bool    { return c == ')' }
func (dsl *dslCollection) isString(c byte) bool     { return c == '"' }
func (dsl *dslCollection) isAssign(c byte) bool     { return c == ':' }
func (dsl *dslCollection) isNamedArg(c byte) bool   { return c == '=' }
func (dsl *dslCollection) isTerminator(c byte) bool { return c == ';' }
func (dsl *dslCollection) isEscape(c byte) bool     { return c == '\\' }
func (dsl *dslCollection) isEmpty(s string) bool    { return len(s) == 0 }
func (dsl *dslCollection) isNewline(s string) bool  { return s == "\n" || s == "\r\n" }
func (dsl *dslCollection) isComment(c byte) bool    { return c == '#' }
func (dsl *dslCollection) isArgRef(c byte) bool     { return c == '$' }
