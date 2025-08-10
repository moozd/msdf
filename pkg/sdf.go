package msdf

import (
	"fmt"
	"image/color"
	"os"

	"golang.org/x/image/font/sfnt"
)

type Msdf struct {
	font *sfnt.Font
	cfg  *Config
}

type Config struct {
	LineHeight, Advance int
	Debug               bool
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
	tex := NewGlyph(m.cfg)
	dbg := NewGlyph(m.cfg)

	edges, scaler, _ := m.getEdges(r)
	contours := edges.getContours()

	for y := range m.cfg.LineHeight {
		for x := range m.cfg.Advance {

			p := scaler.p2g(x, y)
			flipY := m.cfg.LineHeight - 1 - y

			r := contours.getSignedDistnace(RED, p)
			g := contours.getSignedDistnace(GREEN, p)
			b := contours.getSignedDistnace(BLUE, p)

			tex.Image().Set(x, flipY, color.RGBA{
				normalize(r),
				normalize(g),
				normalize(b),
				255,
			})

		}
	}

	if m.cfg.Debug {
		for _, edge := range edges {
			edge.Curve.Debug(dbg, edge.Color.RGB(), scaler)
		}
		dbg.Save(fmt.Sprintf("assets/%c_debug.png", r))

	}

	return tex
}

func normalize(c float64) uint8 {
	// Standard MSDF encoding: negative distances (inside) = high values (white)
	// Positive distances (outside) = low values (dark)
	normalized := 128.0 - c*64.0

	if normalized < 0 {
		return 0
	}
	if normalized > 255 {
		return 255
	}
	return uint8(normalized)
}
