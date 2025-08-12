package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type Vector struct {
	X, Y float64
}

func vec(x0, y0, x1, y1 float64) *Vector {
	return &Vector{
		X: x1 - x0,
		Y: y1 - y0,
	}
}

func (v *Vector) dot(b *Vector) float64 {
	return v.X*b.X + v.Y*b.Y
}

func (v *Vector) cross(b *Vector) float64 {
	return v.X*b.Y - v.Y*b.X
}

func (v *Vector) dist() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) normalize() *Vector {
	l := v.dist()
	return &Vector{
		X: v.X / l,
		Y: v.Y / l,
	}

}

func (v *Vector) fixed() fixed.Point26_6 {
	return fixed.Point26_6{
		X: pack_i26_6(v.X),
		Y: pack_i26_6(v.Y),
	}
}
