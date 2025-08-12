package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type CurveDef interface {
	Params() any
	GetLowResSample(t fixed.Int26_6) fixed.Point26_6
}

type Curve struct {
	def          CurveDef
	Points       []fixed.Point26_6
	DirectionVec *Vector
}

func NewCurve(def CurveDef) *Curve {
	b := &Curve{
		def: def,
	}

	b.doLowResolutionSampling()

	return b
}

func (c *Curve) IsConnected(c2 *Curve) bool {

	N := len(c.Points)
	a := c.Points[N-1]
	s := c2.Points[0]

	return (s == a)
}

func (c *Curve) doLowResolutionSampling() {

	for i := range 65 {
		t := fixed.Int26_6(i)
		p := c.def.GetLowResSample(t)
		c.Points = append(c.Points, p)
	}

	x0, y0 := unpack_p26_6(c.Points[0])
	x1, y1 := unpack_p26_6(c.Points[len(c.Points)-1])
	c.DirectionVec = vec(x0, y0, x1, y1)
}

func (c *Curve) GetPsudoMinimumDistance(xp, yp float64) (float64, float64, float64) {
	switch s := c.def.Params().(type) {
	case Line:

		xp0, yp0 := unpack_p26_6(s.P0)
		xp1, yp1 := unpack_p26_6(s.P1)

		A := vec(xp0, yp0, xp, yp)
		B := vec(xp0, yp0, xp1, yp1)

		t := A.dot(B) / B.dot(B)
		t = clamp(t, 0, 1)

		cx := xp0 + t*(xp1-xp0)
		cy := yp0 + t*(yp1-yp0)
		C := vec(xp, yp, cx, cy)

		return C.dist(), cx, cy

	case QuadraticBezier:
		return 0.0, 0.0, 0.0

	case CubicBezier:

		return 0.0, 0.0, 0.0

	}

	return 0.0, 0.0, 0.0
}

func (c1 *Curve) IsCorner(c2 *Curve, winding ClockDirection, threshold float64) (bool, float64) {
	v1 := c1.DirectionVec.normalize()
	v2 := c2.DirectionVec.normalize()

	dp := v1.dot(v2) / (v1.dist() * v2.dist())
	angle := math.Acos(dp)

	cs := v1.cross(v2)
	isInterior := (winding == CCW && cs > 0) || (winding == CW && cs < 0)
	if !isInterior {
		angle = 2*math.Pi - angle
	}

	if winding == CCW {
		angle = angle - math.Pi
	}

	deg := angle * 180 / math.Pi
	return deg < threshold, deg

}

func (c *Curve) FindMinDistance(p0 fixed.Point26_6) (float64, fixed.Point26_6) {
	r := math.MaxFloat64
	var p fixed.Point26_6

	for _, p1 := range c.Points {
		d := dist(p0, p1)
		if d < r {
			r = d
			p = p1
		}
	}

	return r, p
}

// ------------------

type CubicBezier struct {
	P0, P1, P2, P3 fixed.Point26_6
}

func (cb CubicBezier) Params() any {
	return cb
}

func (cb *CubicBezier) GetLowResSample(step fixed.Int26_6) fixed.Point26_6 {
	t := unpack_i26_6(step)
	x0, y0 := unpack_p26_6(cb.P0)
	x1, y1 := unpack_p26_6(cb.P1)
	x2, y2 := unpack_p26_6(cb.P2)
	x3, y3 := unpack_p26_6(cb.P3)

	T0 := math.Pow(1-t, 3)
	T1 := math.Pow(1-t, 2) * t * 3
	T2 := math.Pow(t, 2) * (1 - t) * 3
	T3 := math.Pow(t, 3)

	x := T0*x0 + T1*x1 + T2*x2 + T3*x3
	y := T0*y0 + T1*y1 + T2*y2 + T3*y3

	return pack_p26_6(x, y)
}

// --------------

type QuadraticBezier struct {
	P0, P1, P2 fixed.Point26_6
}

func (qb QuadraticBezier) Params() any {
	return qb
}

func (qb *QuadraticBezier) GetLowResSample(step fixed.Int26_6) fixed.Point26_6 {
	t := unpack_i26_6(step)

	x0, y0 := unpack_p26_6(qb.P0)
	x1, y1 := unpack_p26_6(qb.P1)
	x2, y2 := unpack_p26_6(qb.P2)

	T0 := math.Pow(1-t, 2)
	T1 := (1 - t) * t * 2
	T2 := math.Pow(t, 2)

	x := T0*x0 + T1*x1 + T2*x2
	y := T0*y0 + T1*y1 + T2*y2

	return pack_p26_6(x, y)
}

// --------------------

type Line struct {
	P0, P1 fixed.Point26_6
}

func (qb Line) Params() any {
	return qb
}

func (qb *Line) GetLowResSample(step fixed.Int26_6) fixed.Point26_6 {
	t := unpack_i26_6(step)
	x0, y0 := unpack_p26_6(qb.P0)
	x1, y1 := unpack_p26_6(qb.P1)

	x := x0 + t*(x1-x0)
	y := y0 + t*(y1-y0)

	return pack_p26_6(x, y)
}
