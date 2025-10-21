package main

import (
	"fmt"

	"github.com/toxyl/math"
)

type Vector struct {
	X float64
	Y float64
	Z float64
}

func (v *Vector) String() string {
	return fmt.Sprintf("V(%.1f %.1f %.1f)", v.X, v.Y, v.Z)
}

func (v *Vector) Delta(v2 *Vector) *Vector {
	return &Vector{
		X: math.Max(v.X, v2.X) - math.Min(v.X, v2.X),
		Y: math.Max(v.Y, v2.Y) - math.Min(v.Y, v2.Y),
		Z: math.Max(v.Z, v2.Z) - math.Min(v.Z, v2.Z),
	}
}

func (v *Vector) Translate(x, y, z float64) *Vector {
	v.X += x
	v.Y += y
	v.Z += z
	return v
}

func (v *Vector) Norm(x, y, z float64) *Vector {
	v.X /= x
	v.Y /= y
	v.Z /= z
	return v
}

func (v *Vector) Denorm(x, y, z float64) *Vector {
	v.X *= x
	v.Y *= y
	v.Z *= z
	return v
}
