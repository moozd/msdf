package msdf

import (
	"image"
	"image/color"
	"os"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Msdf struct {
	font   *sfnt.Font
	config *Config
}

type Config struct {
	LineHeight, Advance int
}

var colors = []EdgeColor{
	RED | GREEN,
	RED | BLUE,
	BLUE | GREEN,
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
		config: cfg,
		font:   fnt,
	}

	return msdf, nil
}

func (m *Msdf) Get(r rune) *image.RGBA {

	edges, s, _ := m.getEdges(r)

	for i, edge := range edges {
		edge.Paint(colors[i%len(colors)])
	}

	tex := image.NewRGBA(image.Rect(0, 0, m.config.Advance, m.config.LineHeight))

	for y := range m.config.LineHeight {
		for x := range m.config.Advance {

			p := s.scale(x, y)

			r := edges.getSignedDistnace(RED, p)
			g := edges.getSignedDistnace(GREEN, p)
			b := edges.getSignedDistnace(BLUE, p)

			tex.Set(x, y, color.RGBA{
				normlize(r),
				normlize(g),
				normlize(b),
				255,
			})

		}
	}

	return tex
}

func normlize(c fixed.Int26_6) uint8 {
	// Fixed-point MSDF normalization
	// Center at 128 (edge), scale for good contrast
	// c is in 26.6 format (value * 64)
	normalized := 128 + (int(c) >> 1) // Divide by 2 for scaling (was >>6 then *32, now >>1)

	if normalized < 0 {
		return 0
	}
	if normalized > 255 {
		return 255
	}
	return uint8(normalized)
}
