package msdf

import (
	"image"

	"golang.org/x/image/math/fixed"
)

type Metrics struct {
	bounds fixed.Rectangle26_6
	adv    fixed.Int26_6
	config *Config
}

func (m *Msdf) getMetrics(r rune) (*Metrics, error) {

	_, bounds, adv, err := m.getVector(r)
	if err != nil {
		return nil, err
	}
	metrics := newMetrics(m.cfg, bounds, adv)
	return metrics, nil
}

func newMetrics(cfg *Config, bounds fixed.Rectangle26_6, adv fixed.Int26_6) *Metrics {
	m := &Metrics{}

	// Store original glyph bounds without padding
	m.bounds = bounds
	m.config = cfg
	m.adv = adv

	return m
}

func (e *Metrics) GetAdvance() float64 {
	return unpack_i26_6(e.adv)
}

func (e *Metrics) GetPlaneBounds() fixed.Rectangle26_6 {
	return e.bounds
}

func (e *Metrics) GetBounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(e.bounds.Min.X.Floor()),
			Y: int(e.bounds.Min.Y.Floor()),
		},
		Max: image.Point{
			X: int(e.bounds.Max.X.Ceil()),
			Y: int(e.bounds.Max.Y.Ceil()),
		},
	}
}

func (e *Metrics) GetRange() (float64, float64) {
	// Return the original glyph dimensions
	x0, y0 := unpack_p26_6(e.bounds.Min)
	x1, y1 := unpack_p26_6(e.bounds.Max)
	rangeX := x1 - x0
	rangeY := y1 - y0
	return rangeX, rangeY
}

func (e *Metrics) ToFloat(x, y int) (float64, float64) {
	// Simple linear mapping from texture space to glyph coordinate space
	rangeX, rangeY := e.GetRange()

	// Normalize texture coordinates [0, width] -> [0, 1] -> glyph bounds
	normalizedX := float64(x) / float64(e.config.width)
	normalizedY := float64(y) / float64(e.config.height)

	// Map to glyph coordinate space
	fx := unpack_i26_6(e.bounds.Min.X) + normalizedX*rangeX
	fy := unpack_i26_6(e.bounds.Max.Y) - normalizedY*rangeY

	return fx, fy
}

func (e *Metrics) Scale(p fixed.Point26_6, bounds image.Rectangle, padding int) (int, int) {
	rangeX, rangeY := e.GetRange()

	// Convert from glyph coords back to texture pixel coords
	normalizedX := unpack_i26_6(p.X-e.bounds.Min.X) / rangeX
	normalizedY := unpack_i26_6(p.Y-e.bounds.Min.Y) / rangeY

	w := bounds.Max.X - bounds.Min.X - 2*padding
	h := bounds.Max.Y - bounds.Min.Y - 2*padding

	px := int(normalizedX*float64(w)) + bounds.Min.X + padding
	py := int(normalizedY*float64(h)) + bounds.Min.Y + padding

	return px, py
}
