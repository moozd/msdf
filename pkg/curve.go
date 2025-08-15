package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type CurveSampler interface {
	PointAt(t float64) Point
	TangentAt(t float64) Point
}

type Curve interface {
	CurveSampler
	IsCorner(c2 Curve, winding ClockDirection, threshold float64) (bool, float64)
	IsConnected(c Curve) bool
	GetSignedArea() float64
	GetLowResPoints() []fixed.Point26_6
	GetDirectionVector() *Vector
}

type baseCurve struct {
	points       []fixed.Point26_6
	directionVec *Vector
}

func (c *baseCurve) IsConnected(c2 Curve) bool {

	N := len(c.points)
	a := c.points[N-1]
	s := c2.GetLowResPoints()[0]

	return (s == a)
}

func (c *baseCurve) doLowResSampling(sampler CurveSampler) {

	for i := range 65 {
		t := float64(i) / 64.0
		p := sampler.PointAt(t)
		c.points = append(c.points, p.fixed())
	}

	x0, y0 := unpack_p26_6(c.points[0])
	x1, y1 := unpack_p26_6(c.points[len(c.points)-1])
	c.directionVec = vec(x0, y0, x1, y1)
}

func (c *baseCurve) GetSignedArea() float64 {
	sum := 0.0
	points := c.points
	N := len(points)

	for i := 1; i < N-2; i += 1 {
		A := points[i]
		B := points[i-1]
		C := points[i+1]

		AB := vec_p26_6(A, B)
		AC := vec_p26_6(A, C)
		sum += AB.Cross(AC)
	}

	return sum
}

func (c baseCurve) GetLowResPoints() []fixed.Point26_6 {
	return c.points
}

func (c baseCurve) GetDirectionVector() *Vector {
	return c.directionVec
}

func (c1 *baseCurve) IsCorner(c2 Curve, winding ClockDirection, threshold float64) (bool, float64) {
	v1 := c1.GetDirectionVector().Normalize()
	v2 := c2.GetDirectionVector().Normalize()

	// Signed angle in radians
	cross := v1.Cross(v2)
	dot := v1.Dot(v2)
	angle := math.Atan2(cross, dot) // signed turn from v1 to v2

	// Adjust based on winding
	if winding == CW {
		angle = -angle
	}

	deg := angle * 180 / math.Pi
	return math.Abs(deg) < threshold, deg

}

// ------------------

type CubicBezier struct {
	P0, P1, P2, P3 fixed.Point26_6
	baseCurve
}

func NewCubicBezier(p0, p1, p2, p3 fixed.Point26_6) *CubicBezier {
	cb := &CubicBezier{
		P0:        p0,
		P1:        p1,
		P2:        p2,
		P3:        p3,
		baseCurve: baseCurve{},
	}

	cb.doLowResSampling(cb)

	return cb
}

func (qc *CubicBezier) TangentAt(t float64) Point {

	x0, y0 := unpack_p26_6(qc.P0)
	x1, y1 := unpack_p26_6(qc.P1)
	x2, y2 := unpack_p26_6(qc.P2)
	x3, y3 := unpack_p26_6(qc.P3)

	x := 3*(x1-x0) + 6*t*(x2-2*x1+x0) + 3*math.Pow(t, 2)*(x3-3*x2+3*x1-x0)
	y := 3*(y1-y0) + 6*t*(y2-2*y1+y0) + 3*math.Pow(t, 2)*(y3-3*y2+3*y1-y0)

	return Point{
		X: x,
		Y: y,
	}
}

func (qc *CubicBezier) PointAt(t float64) Point {
	x0, y0 := unpack_p26_6(qc.P0)
	x1, y1 := unpack_p26_6(qc.P1)
	x2, y2 := unpack_p26_6(qc.P2)
	x3, y3 := unpack_p26_6(qc.P3)

	T0 := math.Pow(1-t, 3)
	T1 := math.Pow(1-t, 2) * t * 3
	T2 := math.Pow(t, 2) * (1 - t) * 3
	T3 := math.Pow(t, 3)

	x := T0*x0 + T1*x1 + T2*x2 + T3*x3
	y := T0*y0 + T1*y1 + T2*y2 + T3*y3

	return Point{
		X: x,
		Y: y,
	}
}

// --------------

type QuadraticBezier struct {
	P0, P1, P2 fixed.Point26_6
	baseCurve
}

func NewQuadraticBezier(p0, p1, p2 fixed.Point26_6) *QuadraticBezier {
	qb := &QuadraticBezier{
		P0:        p0,
		P1:        p1,
		P2:        p2,
		baseCurve: baseCurve{},
	}

	qb.doLowResSampling(qb)

	return qb
}

func (qb *QuadraticBezier) TangentAt(t float64) Point {

	x0, y0 := unpack_p26_6(qb.P0)
	x1, y1 := unpack_p26_6(qb.P1)
	x2, y2 := unpack_p26_6(qb.P2)

	x := 2*(x1-x0) + 2*t*(x2-2*x1+x0)
	y := 2*(y1-y0) + 2*t*(y2-2*y1+y0)

	return Point{
		X: x,
		Y: y,
	}
}

func (qb *QuadraticBezier) PointAt(t float64) Point {

	x0, y0 := unpack_p26_6(qb.P0)
	x1, y1 := unpack_p26_6(qb.P1)
	x2, y2 := unpack_p26_6(qb.P2)

	T0 := math.Pow(1-t, 2)
	T1 := (1 - t) * t * 2
	T2 := math.Pow(t, 2)

	x := T0*x0 + T1*x1 + T2*x2
	y := T0*y0 + T1*y1 + T2*y2

	return Point{
		X: x,
		Y: y,
	}
}

// --------------------

type Line struct {
	P0, P1 fixed.Point26_6
	baseCurve
}

func NewLine(p0, p1 fixed.Point26_6) *Line {
	ln := &Line{
		P0:        p0,
		P1:        p1,
		baseCurve: baseCurve{},
	}

	ln.doLowResSampling(ln)

	return ln
}

func (ln *Line) PointAt(t float64) Point {
	x0, y0 := unpack_p26_6(ln.P0)
	x1, y1 := unpack_p26_6(ln.P1)

	x := x0 + t*(x1-x0)
	y := y0 + t*(y1-y0)

	return Point{
		X: x,
		Y: y,
	}
}

func (l *Line) TangentAt(t float64) Point {
	x0, y0 := unpack_p26_6(l.P0)
	x1, y1 := unpack_p26_6(l.P1)
	return Point{
		X: x1 - x0,
		Y: y1 - y0,
	}
}
