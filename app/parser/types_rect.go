package main

import (
	"fmt"

	"github.com/toxyl/math"
)

type Rect struct {
	P1 *Point
	P2 *Point
}

func (r *Rect) String() string {
	return fmt.Sprintf("R(%f %f %f %f)", r.X1(), r.Y1(), r.X2(), r.Y2())
}

func (r *Rect) X1() float64 { return math.Min(r.P1.X, r.P2.X) }
func (r *Rect) X2() float64 { return math.Max(r.P1.X, r.P2.X) }
func (r *Rect) Y1() float64 { return math.Min(r.P1.Y, r.P2.Y) }
func (r *Rect) Y2() float64 { return math.Max(r.P1.Y, r.P2.Y) }
func (r *Rect) W() float64  { return r.X2() - r.X1() }
func (r *Rect) H() float64  { return r.Y2() - r.Y1() }

func (r *Rect) Delta(r2 *Rect) *Rect {
	return &Rect{
		P1: r.P1.Delta(r2.P1),
		P2: r.P2.Delta(r2.P2),
	}
}

func (r *Rect) Translate(delta Point) *Rect {
	r.P1.Translate(delta.X, delta.Y)
	r.P2.Translate(delta.X, delta.Y)
	return r
}

func (r *Rect) Norm(x, y float64) *Rect {
	r.P1.Norm(x, y)
	r.P2.Norm(x, y)
	return r
}

func (r *Rect) Denorm(x, y float64) *Rect {
	r.P1.Norm(x, y)
	r.P2.Norm(x, y)
	return r
}
