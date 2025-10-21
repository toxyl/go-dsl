package main

import (
	"fmt"
)

type Quad struct {
	P1 *Point
	P2 *Point
	P3 *Point
	P4 *Point
}

func (q *Quad) String() string {
	return fmt.Sprintf("Q(%s %s %s %s)", q.P1.String(), q.P2.String(), q.P3.String(), q.P4.String())
}

func (p *Quad) Delta(p2 *Quad) *Quad {
	return &Quad{
		P1: p.P1.Delta(p2.P1),
		P2: p.P2.Delta(p2.P2),
		P3: p.P3.Delta(p2.P3),
		P4: p.P4.Delta(p2.P4),
	}
}

func (q *Quad) Translate(delta Point) *Quad {
	q.P1.Translate(delta.X, delta.Y)
	q.P2.Translate(delta.X, delta.Y)
	q.P3.Translate(delta.X, delta.Y)
	q.P4.Translate(delta.X, delta.Y)
	return q
}

func (q *Quad) Norm(x, y float64) *Quad {
	q.P1.Norm(x, y)
	q.P2.Norm(x, y)
	q.P3.Norm(x, y)
	q.P4.Norm(x, y)
	return q
}

func (q *Quad) Denorm(x, y float64) *Quad {
	q.P1.Norm(x, y)
	q.P2.Norm(x, y)
	q.P3.Norm(x, y)
	q.P4.Norm(x, y)
	return q
}
