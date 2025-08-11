package msdf

import (
	"math"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Metrics struct {
	bounds fixed.Rectangle26_6
	config *Config
}

func NewMetrics(cfg *Config, segments sfnt.Segments) *Metrics {

	m := &Metrics{}
	bounds := fixed.Rectangle26_6{
		Min: fixed.Point26_6{X: fixed.Int26_6(1 << 20), Y: fixed.Int26_6(1 << 20)},
		Max: fixed.Point26_6{X: fixed.Int26_6(-1 << 20), Y: fixed.Int26_6(-1 << 20)},
	}

	// First pass: calculate bounds
	for _, segment := range segments {
		for _, arg := range segment.Args {
			if arg.X < bounds.Min.X {
				bounds.Min.X = arg.X
			}
			if arg.Y < bounds.Min.Y {
				bounds.Min.Y = arg.Y
			}
			if arg.X > bounds.Max.X {
				bounds.Max.X = arg.X
			}
			if arg.Y > bounds.Max.Y {
				bounds.Max.Y = arg.Y
			}
		}
	}

	padding := pack_i26_6(cfg.Padding)
	bounds.Min.X -= padding
	bounds.Min.Y -= padding
	bounds.Max.X += padding
	bounds.Max.Y += padding

	m.bounds = bounds
	m.config = cfg

	return m
}

func (e *Metrics) GetRange() (fixed.Int26_6, fixed.Int26_6) {

	rangeX := e.bounds.Max.X - e.bounds.Min.X
	rangeY := e.bounds.Max.Y - e.bounds.Min.Y

	return rangeX, rangeY
}

func (e *Metrics) Normalize(t float64) float64 {

	x, y := e.config.Advance-1, e.config.LineHeight-1
	m := dist(e.ToNative(-x, -y), e.ToNative(x, y))

	n := (t + m) / (2 * m)
	n = math.Max(0.0, math.Min(1.0, n))

	return n

}

func (e *Metrics) ToNative(x, y int) fixed.Point26_6 {
	rangeX, rangeY := e.GetRange()

	fx := (fixed.I(x)*rangeX)/fixed.I(e.config.Advance) + e.bounds.Min.X
	fy := e.bounds.Max.Y - (fixed.I(y)*rangeY)/fixed.I(e.config.LineHeight)

	return fixed.Point26_6{
		X: fx,
		Y: fy,
	}
}

func (e *Metrics) ToPixel(p fixed.Point26_6) (int, int) {
	rangeX, rangeY := e.GetRange()

	// Convert back from glyph coords to pixel coords
	px := ((p.X - e.bounds.Min.X) * fixed.I(e.config.Advance)) / rangeX
	py := ((p.Y - e.bounds.Min.Y) * fixed.I(e.config.LineHeight)) / rangeY

	return px.Round(), py.Round()
}
