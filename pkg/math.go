package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

type Shape interface {
	RayHits(p fixed.Point26_6) int
	GetDistance(p fixed.Point26_6) fixed.Int26_6
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

func (s *CubicBezierShape) RayHits(p fixed.Point26_6) int {

	xp, yp := unpack_p26_6(p)

	for t := 0.0; t <= 1.0; t += 0.01 {

		c := s.getCubicBezierPoint(t)
		xc, yc := unpack_p26_6(c)

		if math.Abs(yc-yp) < 1.0 {
			if xc > xp {
				return 1
			}
		}
	}
	return 0
}

func (s *QuadraticBezierShape) RayHits(p fixed.Point26_6) int {
	xp, yp := unpack_p26_6(p)

	for t := 0.0; t <= 1.0; t += 0.01 {

		c := s.getQuadraticBezierPoint(t)
		xc, yc := unpack_p26_6(c)

		if math.Abs(yc-yp) < 1.0 {
			if xc > xp {
				return 1
			}
		}
	}
	return 0
}

func (s *LineShape) RayHits(p fixed.Point26_6) int {
	px, py := unpack_p26_6(p)
	x1, y1 := unpack_p26_6(s.P0)
	x2, y2 := unpack_p26_6(s.P1)

	if (y1 > py) != (y2 > py) {
		intersectX := x1 + (py-y1)/(y2-y1)*(x2-x1)
		if intersectX > px {
			return 1
		}
	}
	return 0
}

func (q *CubicBezierShape) GetDistance(p fixed.Point26_6) fixed.Int26_6 {
	res := fixed.Int26_6(math.MaxInt32)
	for i := range 65 {
		t := float64(i) / 64.0
		p0 := q.getCubicBezierPoint(t)
		d := getDistance(p, p0)

		if d < res {
			res = d
		}
	}
	return res
}

func (q *QuadraticBezierShape) GetDistance(p fixed.Point26_6) fixed.Int26_6 {
	res := fixed.Int26_6(math.MaxInt32)
	for i := range 65 {
		t := float64(i) / 64.0
		p0 := q.getQuadraticBezierPoint(t)
		d := getDistance(p, p0)

		if d < res {
			res = d
		}
	}
	return res
}

func (l *LineShape) GetDistance(p fixed.Point26_6) fixed.Int26_6 {
	x0, y0 := unpack_p26_6(p)
	x1, y1 := unpack_p26_6(l.P0)
	x2, y2 := unpack_p26_6(l.P1)

	n := math.Abs((y2-y1)*x0 - (x2-x1)*y0 + x2*y1 - y2*x1)
	d := math.Sqrt(math.Pow(y2-y1, 2) + math.Pow(x2-x1, 2))

	if d == 0 {
		return 0
	}

	res := n / d
	return pack_i26_6(res)
}

func (c *CubicBezierShape) getCubicBezierPoint(t float64) fixed.Point26_6 {
	x0, y0 := unpack_p26_6(c.P0)
	x1, y1 := unpack_p26_6(c.P1)
	x2, y2 := unpack_p26_6(c.P2)
	x3, y3 := unpack_p26_6(c.P3)

	T0 := math.Pow(1-t, 3)
	T1 := math.Pow(1-t, 2) * t * 3
	T2 := math.Pow(t, 2) * (1 - t) * 3
	T3 := math.Pow(t, 3)

	x := T0*x0 + T1*x1 + T2*x2 + T3*x3
	y := T0*y0 + T1*y1 + T2*y2 + T3*y3

	return pack_p26_6(x, y)
}

func (q *QuadraticBezierShape) getQuadraticBezierPoint(t float64) fixed.Point26_6 {
	x0, y0 := unpack_p26_6(q.P0)
	x1, y1 := unpack_p26_6(q.P1)
	x2, y2 := unpack_p26_6(q.P2)

	T0 := math.Pow(1-t, 2)
	T1 := 2 * t * (1 - t)
	T2 := math.Pow(t, 2)

	x := T0*x0 + T1*x1 + T2*x2
	y := T0*y0 + T1*y1 + T2*y2

	return pack_p26_6(x, y)
}

func getDistance(p0, p1 fixed.Point26_6) fixed.Int26_6 {
	x0, y0 := unpack_p26_6(p0)
	x1, y1 := unpack_p26_6(p1)

	res := math.Sqrt(math.Pow(y1-y0, 2) + math.Pow(x1-x0, 2))

	return pack_i26_6(res)
}

func unpack_p26_6(p fixed.Point26_6) (float64, float64) {
	return float64(p.X) / 64.0, float64(p.Y) / 64.0
}

func pack_p26_6(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
}

func pack_i26_6(f float64) fixed.Int26_6 {
	return fixed.Int26_6(f * 64)
}

func pow26_6(x fixed.Int26_6, a int) fixed.Int26_6 {
	if a == 0 {
		return fixed.I(1)
	}
	if a == 1 {
		return x
	}
	if a < 0 {
		return fixed.I(1) / pow26_6(x, -a)
	}

	result := fixed.Int26_6(1)
	base := x
	exp := a

	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) >> 6
		}
		base = (base * base) >> 6
		exp >>= 1
	}

	return result
}
