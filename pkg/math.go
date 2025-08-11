package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

func dist(p0, p1 fixed.Point26_6) float64 {
	x0, y0 := unpack_p26_6(p0)
	x1, y1 := unpack_p26_6(p1)

	res := math.Sqrt(math.Pow(y1-y0, 2) + math.Pow(x1-x0, 2))

	return res
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func sign(x float64) float64 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}

	return 0
}

func signedArea(points []fixed.Point26_6) fixed.Int26_6 {

	res := fixed.I(0)
	for i, p0 := range points {
		j := (i + 1) % len(points)
		p1 := points[j]

		res += (p1.X - p0.X) * (-p1.Y + -p0.Y)
	}

	return res
}

func angle_abc(p0, p1, p2 fixed.Point26_6) float64 {
	x0, y0 := unpack_p26_6(p0)
	x1, y1 := unpack_p26_6(p1)
	x2, y2 := unpack_p26_6(p2)

	ax, ay, bx, by := (x0 - x1), (y0 - y1), (x2 - x1), (y2 - y1)
	dp := ax*bx + ay*by
	ma := math.Sqrt(math.Pow(ax, 2) + math.Pow(ay, 2))
	mb := math.Sqrt(math.Pow(bx, 2) + math.Pow(by, 2))

	res := dp / (ma * mb)
	ang := math.Acos(res) * (180 * math.Pi)

	// fmt.Printf("angle: A(%f,%f) B(%f,%f) C(%f,%f)\n", x0, y0, x1, y1, x2, y2)
	// fmt.Printf("angle: %f deg\n", ang)

	return ang
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

func unpack_i26_6(f fixed.Int26_6) float64 {
	return float64(f) / 64.0
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
