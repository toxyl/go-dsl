package main

import (
	"fmt"
)

type NGon struct {
	Points []*Point
}

func (n *NGon) String() string {
	return fmt.Sprintf("N(%d points)", len(n.Points))
}

func (n *NGon) Delta(n2 NGon) *NGon {
	for i, p := range n.Points {
		n.Points[i] = p.Delta(n2.Points[i])
	}
	return n
}

func (n *NGon) Translate(x, y float64) *NGon {
	for _, p := range n.Points {
		p.Translate(x, y)
	}
	return n
}

func (n *NGon) Norm(x, y float64) *NGon {
	for _, p := range n.Points {
		p.Norm(x, y)
	}
	return n
}

func (n *NGon) Denorm(x, y float64) *NGon {
	for _, p := range n.Points {
		p.Denorm(x, y)
	}
	return n
}
