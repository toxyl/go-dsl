package main

import (
	"fmt"
	"strings"
)

const CONTEXT_CHARS = 25

// formatErrorWithPosition wraps an error with line/column position information and character-based context.
// It extracts ~CONTEXT_CHARS characters before and after the error position and formats the error message
// with a visual indicator (^) pointing to the error location.
func formatErrorWithPosition(err error, source string, line, col int) error {
	if err == nil {
		return nil
	}

	// Convert line/column to character position in source
	charPos := lineColToCharPos(source, line, col)
	if charPos < 0 || charPos >= len(source) {
		return fmt.Errorf("[%d:%d] %v", line, col, err)
	}

	// Extract context (~CONTEXT_CHARS chars before and after)
	contextStart := charPos - CONTEXT_CHARS
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := charPos + CONTEXT_CHARS
	if contextEnd > len(source) {
		contextEnd = len(source)
	}

	context := source[contextStart:contextEnd]
	indicatorPos := charPos - contextStart

	// Replace newlines in context with spaces for display
	displayContext := strings.ReplaceAll(context, "\n", " ")
	displayContext = strings.ReplaceAll(displayContext, "\r", " ")

	// Build the indicator line
	indicator := make([]byte, len(displayContext))
	for i := range indicator {
		if i == indicatorPos && indicatorPos < len(displayContext) {
			indicator[i] = '^'
		} else {
			indicator[i] = ' '
		}
	}

	return fmt.Errorf("[%d:%d] %v, check around:\n`%s`\n %s", line, col, err, displayContext, string(indicator))
}

// lineColToCharPos converts line/column (1-based) to character position (0-based) in source.
func lineColToCharPos(source string, line, col int) int {
	if line < 1 || col < 1 {
		return -1
	}

	currentLine := 1
	currentCol := 1

	for i := 0; i < len(source); i++ {
		if currentLine == line && currentCol == col {
			return i
		}

		switch source[i] {
		case '\n':
			currentLine++
			currentCol = 1
		case '\r':
			// Handle \r\n sequence
			if i+1 < len(source) && source[i+1] == '\n' {
				currentLine++
				currentCol = 1
				i++ // Skip the \n
			} else {
				currentCol = 1
			}
		default:
			currentCol++
		}
	}

	// If we reached the end and the position is at the end
	if currentLine == line && currentCol == col {
		return len(source)
	}

	return -1
}
