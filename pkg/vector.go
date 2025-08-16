package msdf

import (
	"fmt"
	"math"

	"golang.org/x/image/math/fixed"
)

type Vector struct {
	X, Y float64
}

func vec() *Vector {
	return &Vector{}
}

func (v *Vector) fromXY(x0, y0, x1, y1 float64) *Vector {
	v.X = x1 - x0
	v.Y = y1 - y0

	return v
}

func (v *Vector) fromAB(a, b Point) *Vector {
	v.fromXY(a.X, a.Y, b.X, b.Y)
	return v
}

func (v *Vector) fromP(b Point) *Vector {
	v.fromXY(0, 0, b.X, b.Y)
	return v
}

func (v *Vector) fromP26_6(a, b fixed.Point26_6) *Vector {
	xa, ya := unpack_p26_6(a)
	xb, yb := unpack_p26_6(b)
	v.fromXY(xa, ya, xb, yb)
	return v
}

func (v *Vector) Dot(b *Vector) float64 {
	return v.X*b.X + v.Y*b.Y
}

func (v *Vector) Sub(b *Vector) *Vector {
	return &Vector{
		X: v.X - b.X,
		Y: v.Y - b.Y,
	}
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
