package msdf

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/image/font/sfnt"
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

			if x < 5 && y < 5 { // Debug first few pixels
				fmt.Printf("Pixel (%d,%d): r=%.3f g=%.3f b=%.3f\n", x, y, r, g, b)
			}

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

func normlize(c float64) uint8 {
	// More aggressive mapping for visibility
	if c < 0 {
		return 0 // Inside = black
	} else {
		return 255 // Outside = white  
	}
}
