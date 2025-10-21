package main

import (
	"fmt"
	"image/color"
)

type LineStyle struct {
	Thickness float64
	Color     *color.RGBA64
}

func (l *LineStyle) String() string {
	return fmt.Sprintf("LS(%.1f %v)", l.Thickness, l.Color)
}

type FillStyle struct {
	Color *color.RGBA64
}

func (f *FillStyle) String() string {
	return fmt.Sprintf("FS(%v)", f.Color)
}

type TextStyle struct {
	Family string
	Size   float64
	Color  *color.RGBA64
}

func (f *TextStyle) String() string {
	return fmt.Sprintf("TS(%v)", f.Color)
}
