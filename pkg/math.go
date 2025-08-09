package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

func getDistance(p0, p1 fixed.Point26_6) float64 {
	x0, y0 := unpack_p26_6(p0)
	x1, y1 := unpack_p26_6(p1)

	res := math.Sqrt(math.Pow(y1-y0, 2) + math.Pow(x1-x0, 2))

	return res
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
