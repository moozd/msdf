package msdf

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
		config: cfg,
		font:   fnt,
	}

	return msdf, nil
}

func (m *Msdf) Get(r rune) *image.RGBA {

	edges, scaler, _ := m.getEdges(r)

	var corners [][]*Edge

	corners = make([][]*Edge, len(edges))
	for i := range edges {
		for j := range edges {
			b := edges[i].Curve.IsAttached(edges[j].Curve)
			if b {
				corners[i] = append(corners[i], &edges[j])
			}
		}
		corners[i] = append(corners[i], &edges[i])
	}

	mp := make(map[*Edge]bool)

	for i, c := range corners {
		for j := range c {

			// if mp[corners[i][j]] {
			// 	continue
			// }
			cl := colors[j%len(colors)]

			corners[i][j].Color = cl
			mp[corners[i][j]] = true
		}
	}

	for i := range edges {
		fmt.Println(corners[i])
	}

	tex := image.NewRGBA(image.Rect(0, 0, m.config.Advance, m.config.LineHeight))
	bg := &image.Uniform{color.RGBA{0, 0, 0, 255}}
	draw.Draw(tex, tex.Bounds(), bg, image.Point{}, draw.Src)

	// for y := range m.config.LineHeight {
	// 	for x := range m.config.Advance {
	//
	// 		p := scaler.p2g(x, y)
	//
	// 		r := edges.getSignedDistnace(RED, p)
	// 		g := edges.getSignedDistnace(GREEN, p)
	// 		b := edges.getSignedDistnace(BLUE, p)
	//
	// 		if math.Abs(r-g) > 1 || math.Abs(r-b) > 1 || math.Abs(g-b) > 1 {
	// 			fmt.Printf("Pixel (%d,%d): R=%.2f G=%.2f B=%.2f\n", x, y, r, g, b)
	// 		}
	// 		tex.Set(x, y, color.RGBA{
	// 			normlize(r),
	// 			normlize(g),
	// 			normlize(b),
	// 			255,
	// 		})
	//
	// 	}
	// }

	for _, edge := range edges {
		edge.Curve.Debug(tex, edge.Color.RGB(), scaler)
	}

	return tex
}

func normlize(c float64) uint8 {
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
