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

	metrics, _ := m.getMetrics(r)
	contours, _ := m.getContours(r)

	for i, con := range contours {
		fmt.Printf("con: %d, dir: %v\n", i+1, con.winding)
	}

	for y := range m.cfg.LineHeight {
		for x := range m.cfg.Advance {

			xi, yi := metrics.ToFloat(x, y)
			flippedY := m.cfg.LineHeight - 1 - y

			r := getDistance(contours, RED, xi, yi)
			g := getDistance(contours, GREEN, xi, yi)
			b := getDistance(contours, BLUE, xi, yi)

			tex.Image().Set(x, flippedY, color.RGBA{r, g, b, 255})

		}
	}

	if m.cfg.Debug {
		for _, con := range contours {
			for _, edge := range con.edges {
				edge.Curve.Debug(dbg, edge.Color.RGB(), metrics)
			}
		}
		dbg.Save(fmt.Sprintf("assets/%c_debug.png", r))

	}
	return tex
}

func getDistance(contours []*Contour, c EdgeColor, x, y float64) uint8 {
	var edgeDirectionVec *Vector
	var winding ClockDirection
	var x1, y1 float64

	distance := math.MaxFloat64

	for _, con := range contours {

		winding = con.winding
		for _, edge := range con.edges {
			if !edge.Color.Has(c) {
				continue
			}
			d, xp, yp := edge.Curve.GetPsudoMinimumDistance(x, y)
			if d < distance {
				distance = d
				edgeDirectionVec = edge.Curve.DirectionVec
				x1 = xp
				y1 = yp
			}
		}

	}

	fmt.Printf("distance: %f\n", distance)

	pointVec := vec(x, y, x1, y1)

	side := sign(edgeDirectionVec.cross(pointVec))

	distance = side * float64(winding) * distance

	distanceRange := 0.12

	n := (distance / distanceRange) + 0.5
	n = clamp(n, 0, 1)

	return uint8(n * 255)
}
