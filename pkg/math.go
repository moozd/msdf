package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type Shape interface {
	Intersect(p fixed.Point26_6) int
	GetDistance(p fixed.Point26_6) float64
}

type CubicBezierShape struct {
	P0, P1, P2, P3 fixed.Point26_6
}

type QuadraticBezierShape struct {
	P0, P1, P2 fixed.Point26_6
}

type LineShape struct {
	P0, P1 fixed.Point26_6
}

func (s *CubicBezierShape) Intersect(p fixed.Point26_6) int {

	xp, yp := p26_6__f64(p)

	for i := range 65 {
		t := float64(i) / 64.0

		c := s.getCubicBezierPoint(t)
		xc, yc := p26_6__f64(c)

		if math.Abs(yc-yp) < 1.0 {
			if xc > xp {
				return 1
			}
		}
	}
	return 0
}

func (s *QuadraticBezierShape) Intersect(p fixed.Point26_6) int {
	xp, yp := p26_6__f64(p)

	for i := range 65 {
		t := float64(i) / 64.0

		c := s.getQuadraticBezierPoint(t)
		xc, yc := p26_6__f64(c)

		if math.Abs(yc-yp) < 1.0 {
			if xc > xp {
				return 1
			}
		}
	}
	return 0
}

func (s *LineShape) Intersect(p fixed.Point26_6) int {
	px, py := p26_6__f64(p)
	x1, y1 := p26_6__f64(s.P0)
	x2, y2 := p26_6__f64(s.P1)

	if (y1 > py) != (y2 > py) {
		intersectX := x1 + (py-y1)/(y2-y1)*(x2-x1)
		if intersectX > px {
			return 1
		}
	}
	return 0
}

func (q *CubicBezierShape) GetDistance(p fixed.Point26_6) float64 {
	res := fixed.Int26_6(math.MaxInt32)
	for i := range 65 {
		t := float64(i) / 64.0
		p0 := q.getCubicBezierPoint(t)
		d := getDistance(p, p0)

		if d < res {
			res = d
		}
	}
	return i26_6__f64(res)
}

func (q *QuadraticBezierShape) GetDistance(p fixed.Point26_6) float64 {
	res := fixed.Int26_6(math.MaxInt32)
	for i := range 65 {
		t := float64(i) / 64.0
		p0 := q.getQuadraticBezierPoint(t)
		d := getDistance(p, p0)

		if d < res {
			res = d
		}
	}
	return i26_6__f64(res)
}

func (l *LineShape) GetDistance(p fixed.Point26_6) float64 {
	x0, y0 := p26_6__f64(p)
	x1, y1 := p26_6__f64(l.P0)
	x2, y2 := p26_6__f64(l.P1)

	n := math.Abs((y2-y1)*x0 - (x2-x1)*y0 + x2*y1 - y2*x1)
	d := math.Sqrt(math.Pow(y2-y1, 2) + math.Pow(x2-x1, 2))
	res := n / d

	return res
}

func (c *CubicBezierShape) getCubicBezierPoint(t float64) fixed.Point26_6 {
	x0, y0 := p26_6__f64(c.P0)
	x1, y1 := p26_6__f64(c.P1)
	x2, y2 := p26_6__f64(c.P2)
	x3, y3 := p26_6__f64(c.P3)

	T0 := math.Pow(1-t, 3)
	T1 := 3 * t * math.Pow((1-t), 2)
	T2 := 3 * math.Pow(t, 2) * (1 - t)
	T3 := math.Pow(t, 3)

	x := T0*x0 + T1*x1 + T2*x2 + T3*x3
	y := T0*y0 + T1*y1 + T2*y2 + T3*y3

	return f64__p26_6(x, y)
}

func (q *QuadraticBezierShape) getQuadraticBezierPoint(t float64) fixed.Point26_6 {
	x0, y0 := p26_6__f64(q.P0)
	x1, y1 := p26_6__f64(q.P1)
	x2, y2 := p26_6__f64(q.P2)

	T0 := math.Pow(1-t, 2)
	T1 := 2 * t * (1 - t)
	T2 := math.Pow(t, 2)

	x := T0*x0 + T1*x1 + T2*x2
	y := T0*y0 + T1*y1 + T2*y2

	return f64__p26_6(x, y)
}

func getDistance(p0, p1 fixed.Point26_6) fixed.Int26_6 {
	x0, y0 := p26_6__f64(p0)
	x1, y1 := p26_6__f64(p1)

	res := math.Sqrt(math.Pow(y1-y0, 2) + math.Pow(x1-x0, 2))

	return f64__i26_6(res)
}

func p26_6__f64(p fixed.Point26_6) (float64, float64) {
	return float64(p.X) / 64.0, float64(p.Y) / 64.0
}

func f64__p26_6(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
}

func i26_6__f64(i fixed.Int26_6) float64 {
	return float64(i) / 64.0
}

func f64__i26_6(f float64) fixed.Int26_6 {
	return fixed.Int26_6(f * 64)
}

func f64_sign(f float64) float64 {
	if f >= 0 {
		return 1.0
	}

	return -1.0

}
