package msdf

import (
	"math"

	"golang.org/x/image/math/fixed"
)

const (
	DEG = 180.0 / math.Pi
	RAD = math.Pi / 180.0
)

func sign(s float64) float64 {
	if s > 0 {
		return +1
	}
	if s < 0 {
		return -1
	}
	return 0
}

func lerp(p0, p1 fixed.Point26_6, t float64) fixed.Point26_6 {
	x0, y0 := unpack_p26_6(p0)
	x1, y1 := unpack_p26_6(p1)

	x := x0 + t*(x1-x0)
	y := y0 + t*(y1-y0)

	return fixed.Point26_6{
		X: pack_i26_6(x),
		Y: pack_i26_6(y),
	}

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
