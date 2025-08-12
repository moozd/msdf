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
	height, width int
	Padding       float64
	Debug         bool
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
	metrics, _ := m.getMetrics(r)
	contours, _ := m.getContours(r)

	w, h := metrics.GetRange()

	minSize := 64
	scale := 2.0
	m.cfg.height = max(int(h*scale), minSize) + int(m.cfg.Padding*10)
	m.cfg.width = max(int(w*scale), minSize) + int(m.cfg.Padding*10)

	tex := NewGlyph(m.cfg)
	dbg := NewGlyph(m.cfg)

	fmt.Println(m.cfg)

	for i, con := range contours {
		fmt.Printf("con: %d, dir: %v\n", i+1, con.winding)
	}

	for y := range m.cfg.height {
		for x := range m.cfg.width {

			xi, yi := metrics.ToFloat(x, y)
			flippedY := m.cfg.height - 1 - y
			fmt.Printf("%f,%f\n", xi, yi)

			r := uint8(250) //getDistance(m.cfg, contours, RED, xi, yi)
			g := uint8(250) //getDistance(m.cfg, contours, GREEN, xi, yi)
			b := uint8(25)  //getDistance(m.cfg, contours, BLUE, xi, yi)

			tex.Image().Set(x, flippedY, color.RGBA{r, g, b, 255})

		}

		fmt.Println()
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

func getDistance(cfg *Config, contours []*Contour, c EdgeColor, x, y float64) uint8 {
	var edgeDirectionVec *Vector
	var winding ClockDirection
	var x1, y1 float64
	found := false

	distance := math.MaxFloat64

	for _, con := range contours {

		for _, edge := range con.edges {
			if !edge.Color.Has(c) {
				continue
			}
			d, xp, yp := edge.Curve.GetPsudoMinimumDistance(x, y)
			if d < distance {
				winding = con.winding
				distance = d
				edgeDirectionVec = edge.Curve.DirectionVec
				x1 = xp
				y1 = yp
				found = true
			}
		}

	}

	if !found {
		return 128
	}

	pointVec := vec(x1, y1, x, y)

	side := sign(edgeDirectionVec.cross(pointVec))

	w := 1.0
	if winding == 0 {
		w = -1.0
	}

	distance = side * w * distance
	pixelSize := math.Min(float64(cfg.width), float64(cfg.height))
	distanceRange := (2.0 / pixelSize) * 100

	normalized := (distance / distanceRange) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
