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
	Padding             float64
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

	edges, metrics, _ := m.getEdges(r)
	contours := edges.getContours()

	for y := range m.cfg.LineHeight {
		for x := range m.cfg.Advance {

			p := metrics.ToNative(x, y)
			flipY := m.cfg.LineHeight - 1 - y

			r := contours.getSignedDistnace(metrics, p, RED)
			g := contours.getSignedDistnace(metrics, p, GREEN)
			b := contours.getSignedDistnace(metrics, p, BLUE)

			tex.Image().Set(x, flipY, color.RGBA{
				normalizeColor(x, y, r),
				normalizeColor(x, y, g),
				normalizeColor(x, y, b),
				255,
			})

		}
	}

	if m.cfg.Debug {
		for _, edge := range edges {
			edge.Curve.Debug(dbg, edge.Color.RGB(), metrics)
		}
		dbg.Save(fmt.Sprintf("assets/%c_debug.png", r))

	}

	return tex
}

func normalizeColor(x, y int, c float64) uint8 {
	// Standard MSDF encoding: negative distances (inside) = high values (white)
	// Positive distances (outside) = low values (dark)

	padNoise := float64((x*7+y*13)%100) / 100.0
	noisy := (c + padNoise*0.3) / 1.3
	noisy = math.Max(0, math.Min(1.0, noisy))

	return uint8(noisy * 255)
}
