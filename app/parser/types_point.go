package main

import (
	"fmt"

	"github.com/toxyl/math"
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) String() string {
	return fmt.Sprintf("P(%f %f)", p.X, p.Y)
}

func (p *Point) Delta(p2 *Point) *Point {
	return &Point{
		X: math.Max(p.X, p2.X) - math.Min(p.X, p2.X),
		Y: math.Max(p.Y, p2.Y) - math.Min(p.Y, p2.Y),
	}
}

func (p *Point) Translate(x, y float64) *Point {
	p.X += x
	p.Y += y
	return p
}

func (p *Point) Norm(x, y float64) *Point {
	p.X /= x
	p.Y /= y
	return p
}

func (p *Point) Denorm(x, y float64) *Point {
	p.X *= x
	p.Y *= y
	return p
}
