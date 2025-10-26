package main

import (
	"fmt"
)

type Text struct {
	Text  string
	Style *TextStyle
}

func (t *Text) String() string {
	return fmt.Sprintf("T(\"%s\" %s)", t.Text, t.Style.String())
}
