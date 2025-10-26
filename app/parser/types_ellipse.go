package main

import (
	"fmt"
)

type Ellipse struct {
	Center *Point
	Radius *Point
}

func (e *Ellipse) String() string {
	return fmt.Sprintf("E(%s %s)", e.Center.String(), e.Radius.String())
}

func (n *Ellipse) Delta(n2 Ellipse) *Ellipse {
	return &Ellipse{
		Center: n.Center.Delta(n2.Center),
		Radius: n.Radius.Delta(n2.Radius),
	}
}

func (e *Ellipse) Translate(delta Point) *Ellipse {
	e.Center.Translate(delta.X, delta.Y)
	return e
}

func (e *Ellipse) Norm(x, y float64) *Ellipse {
	e.Center.Norm(x, y)
	return e
}

func (e *Ellipse) Denorm(x, y float64) *Ellipse {
	e.Center.Denorm(x, y)
	return e
}
