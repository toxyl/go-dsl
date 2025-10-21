package main

import (
	"fmt"
	"image"
	"image/color"
)

func (dsl *dslCollection) shellResultImage(img image.Image) string {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	resPrefix := ""
	switch img.(type) {
	case *image.RGBA:
		resPrefix = "image8"
	case *image.RGBA64:
		resPrefix = "image16"
	case *image.NRGBA:
		resPrefix = "image8"
	case *image.NRGBA64:
		resPrefix = "image16"
	}
	return fmt.Sprintf("%s(w: %d, h: %d)", resPrefix, w, h)
}

func (dsl *dslCollection) shellResultColor(col color.Color) string {
	r, g, b, a := col.RGBA()
	switch col.(type) {
	case color.RGBA:
		return fmt.Sprintf("color8(r: %.1f%%, g: %.1f%%, b: %.1f%%, a: %.1f%%)", float64(r)/255.0*100, float64(g)/255.0*100, float64(b)/255.0*100, float64(a)/255.0*100)
	case color.RGBA64:
		return fmt.Sprintf("color16(r: %.1f%%, g: %.1f%%, b: %.1f%%, a: %.1f%%)", float64(r)/65535.0*100, float64(g)/65535.0*100, float64(b)/65535.0*100, float64(a)/65535.0*100)
	case color.NRGBA:
		return fmt.Sprintf("color8(r: %.1f%%, g: %.1f%%, b: %.1f%%, a: %.1f%%)", float64(r)/255.0*100, float64(g)/255.0*100, float64(b)/255.0*100, float64(a)/255.0*100)
	case color.NRGBA64:
		return fmt.Sprintf("color16(r: %.1f%%, g: %.1f%%, b: %.1f%%, a: %.1f%%)", float64(r)/65535.0*100, float64(g)/65535.0*100, float64(b)/65535.0*100, float64(a)/65535.0*100)
	}
	return fmt.Sprintf("color(r: %d, g: %d, b: %d, a: %d)", r, g, b, a)
}

func (dsl *dslCollection) shellResultPoint(p Point) string         { return p.String() }
func (dsl *dslCollection) shellResultRect(r Rect) string           { return r.String() }
func (dsl *dslCollection) shellResultNGon(p NGon) string           { return p.String() }
func (dsl *dslCollection) shellResultTriangle(t Triangle) string   { return t.String() }
func (dsl *dslCollection) shellResultQuad(q Quad) string           { return q.String() }
func (dsl *dslCollection) shellResultEllipse(e Ellipse) string     { return e.String() }
func (dsl *dslCollection) shellResultLineStyle(l LineStyle) string { return l.String() }
func (dsl *dslCollection) shellResultFillStyle(f FillStyle) string { return f.String() }
func (dsl *dslCollection) shellResultTextStyle(t TextStyle) string { return t.String() }
func (dsl *dslCollection) shellResultVector(v Vector) string       { return v.String() }
func (dsl *dslCollection) shellResultText(t Text) string           { return t.String() }

// TODO: NEW TYPES: add additional shellResult* functions
