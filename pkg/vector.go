package msdf

import (
	"fmt"
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

func vec_p26_6(a, b fixed.Point26_6) *Vector {
	xa, ya := unpack_p26_6(a)
	xb, yb := unpack_p26_6(b)
	return vec(xa, ya, xb, yb)
}

func (v *Vector) Dot(b *Vector) float64 {
	return v.X*b.X + v.Y*b.Y
}

func (v *Vector) Cross(b *Vector) float64 {
	return v.X*b.Y - v.Y*b.X
}

func (v *Vector) Distance() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Normalize() *Vector {
	l := v.Distance()
	return &Vector{
		X: v.X / l,
		Y: v.Y / l,
	}

}

func (v *Vector) Fixed() fixed.Point26_6 {
	return fixed.Point26_6{
		X: pack_i26_6(v.X),
		Y: pack_i26_6(v.Y),
	}
}

func (v *Vector) String() string {
	return fmt.Sprintf("vec[(%.3f, %.3f) d= %.3f] ", v.X, v.Y, v.Distance())
}
