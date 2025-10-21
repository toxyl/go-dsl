package main

import (
	"fmt"
)

type Triangle struct {
	P1 *Point
	P2 *Point
	P3 *Point
}

func (t *Triangle) String() string {
	return fmt.Sprintf("T(%s %s %s)", t.P1.String(), t.P2.String(), t.P3.String())
}

func (t *Triangle) Delta(t2 *Triangle) *Triangle {
	return &Triangle{
		P1: t.P1.Delta(t2.P1),
		P2: t.P2.Delta(t2.P2),
		P3: t.P3.Delta(t2.P3),
	}
}

func (t *Triangle) Translate(delta Point) *Triangle {
	t.P1.Translate(delta.X, delta.Y)
	t.P2.Translate(delta.X, delta.Y)
	t.P3.Translate(delta.X, delta.Y)
	return t
}

func (t *Triangle) Norm(x, y float64) *Triangle {
	t.P1.Norm(x, y)
	t.P2.Norm(x, y)
	t.P3.Norm(x, y)
	return t
}

func (t *Triangle) Denorm(x, y float64) *Triangle {
	t.P1.Denorm(x, y)
	t.P2.Denorm(x, y)
	t.P3.Denorm(x, y)
	return t
}
