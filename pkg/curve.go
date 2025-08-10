package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type CurveSampler interface {
	Get(t fixed.Int26_6) fixed.Point26_6
}

type Curve struct {
	sampler CurveSampler
	Points  []fixed.Point26_6
}

func NewCurve(ci CurveSampler) *Curve {
	b := &Curve{
		sampler: ci,
	}

	b.Sample()

	return b
}

func (c *Curve) IsAttached(b *Curve) bool {
	end1 := c.Points[len(c.Points)-1]
	start2, end2 := b.Points[0], b.Points[len(b.Points)-1]

	return (end1 == start2) || (end1 == end2)
}

func (c *Curve) Sample() {

	for i := range 65 {
		t := fixed.Int26_6(i)
		p := c.sampler.Get(t)
		c.Points = append(c.Points, p)
	}
}

func (c *Curve) FindMinDistance(p0 fixed.Point26_6) float64 {
	r := math.MaxFloat64
	for _, p1 := range c.Points {
		r = math.Min(r, getDistance(p0, p1))
	}

	return r
}

func (c *Curve) Cast(p fixed.Point26_6) int {
	winding := 0

	for i := 0; i < len(c.Points)-1; i++ {
		p1 := c.Points[i]
		p2 := c.Points[i+1]

		// Check if horizontal ray from point p crosses the line segment p1->p2
		if (p1.Y > p.Y) != (p2.Y > p.Y) {
			// Calculate X intersection point
			intersectX := p1.X + (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
			if intersectX > p.X {
				// Ray crosses to the right, count it
				if p2.Y > p1.Y {
					winding += 1 // Upward crossing
				} else {
					winding -= 1 // Downward crossing
				}
			}
		}
	}

	return winding
}

// ------------------

type CubicBezier struct {
	P0, P1, P2, P3 fixed.Point26_6
}

func (cb *CubicBezier) Get(step fixed.Int26_6) fixed.Point26_6 {
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

func (qb *QuadraticBezier) Get(step fixed.Int26_6) fixed.Point26_6 {
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

func (qb *Line) Get(step fixed.Int26_6) fixed.Point26_6 {
	t := unpack_i26_6(step)
	x0, y0 := unpack_p26_6(qb.P0)
	x1, y1 := unpack_p26_6(qb.P1)

	x := x0 + t*(x1-x0)
	y := y0 + t*(y1-y0)

	return pack_p26_6(x, y)
}
