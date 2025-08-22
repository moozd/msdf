package msdf

import (
	"fmt"
	"math"

	"golang.org/x/image/math/fixed"
)

type Vector struct {
	X, Y float64
}

type Point struct {
	X, Y float64
}

func (p Point) asVector() *Vector {
	return vec().fromAB(Point{}, p)
}

func (p Point) asFixed() fixed.Point26_6 {
	return packP26_6(p.X, p.Y)
}

func (p Point) Sub(b Point) *Vector {
	return vec().fromAB(p, b)
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

func (v *Vector) fromP26_6(a, b fixed.Point26_6) *Vector {
	xa, ya := unpackP26_6(a)
	xb, yb := unpackP26_6(b)
	v.fromXY(xa, ya, xb, yb)
	return v
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
		X: packI26_6(v.X),
		Y: packI26_6(v.Y),
	}
}

func (v *Vector) String() string {
	return fmt.Sprintf("vec[(%.3f, %.3f) d= %.3f] ", v.X, v.Y, v.Distance())
}
