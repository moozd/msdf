package msdf

import (
	"golang.org/x/image/math/fixed"
)

type Point struct {
	X, Y float64
}

func sign(s float64) float64 {
	if s > 0 {
		return +1
	}
	if s < 0 {
		return -1
	}
	return 0
}

func (p Point) fixed() fixed.Point26_6 {
	return pack_p26_6(p.X, p.Y)
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
