package msdf

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"golang.org/x/image/font/sfnt"
)

type Msdf struct {
	font *sfnt.Font
	cfg  *Config
}

type Config struct {
	LineHeight, Advance int
}

var colors = []EdgeColor{
	RED,
	GREEN,
	BLUE,
}

func New(addr string, cfg *Config) (*Msdf, error) {

	fd, err := os.ReadFile(addr)

	if err != nil {
		return nil, err
	}

	fnt, err := sfnt.Parse(fd)

	if err != nil {
		return nil, err
	}

	msdf := &Msdf{
		cfg:  cfg,
		font: fnt,
	}

	return msdf, nil
}

func (m *Msdf) Get(r rune) *Glyph {

	edges, scaler, _ := m.getEdges(r)

	edges.setupColors()

	tex := NewGlyph(m.cfg)
	dbg := NewGlyph(m.cfg)

	for y := range m.cfg.LineHeight {
		for x := range m.cfg.Advance {

			p := scaler.p2g(x, y)

			r := edges.getSignedDistnace(RED, p)
			g := edges.getSignedDistnace(GREEN, p)
			b := edges.getSignedDistnace(BLUE, p)

			if math.Abs(r-g) > 1 || math.Abs(r-b) > 1 || math.Abs(g-b) > 1 {
				fmt.Printf("PIX (%d,%d): R=%.2f G=%.2f B=%.2f\n", x, y, r, g, b)
			}
			tex.Image().Set(x, y, color.RGBA{
				normalize(r),
				normalize(g),
				normalize(b),
				255,
			})

		}
	}

	for _, edge := range edges {
		edge.Curve.Debug(dbg, edge.Color.RGB(), scaler)
	}

	dbg.Save(fmt.Sprintf("%c_debug.png", r))

	return tex
}

func normalize(c float64) uint8 {
	// Convert distance to range [0, 255] where 128 is the zero-distance point
	// Typical MSDF range is roughly [-4, 4] pixels, so scale accordingly
	normalized := 128.0 + c*16.0 // Scale distance and offset to center at 128
	// fmt.Printf("Distance: %f, Normalized: %f\n", c, normalized)

	if normalized < 0 {
		return 0
	}
	if normalized > 255 {
		return 255
	}
	return uint8(normalized)
}
