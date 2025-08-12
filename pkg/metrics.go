package msdf

import (
	"golang.org/x/image/math/fixed"
)

type Metrics struct {
	bounds fixed.Rectangle26_6
	config *Config
}

func (m *Msdf) getMetrics(r rune) (*Metrics, error) {

	_, bounds, err := m.getVector(r)
	if err != nil {
		return nil, err
	}
	metrics := newMetrics(m.cfg, bounds)
	return metrics, nil
}

func newMetrics(cfg *Config, bounds fixed.Rectangle26_6) *Metrics {
	m := &Metrics{}

	// Store original glyph bounds without padding
	m.bounds = bounds
	m.config = cfg

	return m
}

func (e *Metrics) GetRange() (float64, float64) {
	// Return the original glyph dimensions
	rangeX := e.bounds.Max.X - e.bounds.Min.X
	rangeY := e.bounds.Max.Y - e.bounds.Min.Y

	return unpack_i26_6(rangeX), unpack_i26_6(rangeY)
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

func (e *Metrics) ToPixel(p fixed.Point26_6) (int, int) {
	rangeX, rangeY := e.GetRange()

	// Convert from glyph coords back to texture pixel coords
	normalizedX := float64(p.X-e.bounds.Min.X) / rangeX
	normalizedY := float64(e.bounds.Max.Y-p.Y) / rangeY

	px := int(normalizedX * float64(e.config.width))
	py := int(normalizedY * float64(e.config.height))

	return px, py
}

