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
	Seed           uint
	height, width  int
	Scale          float64
	Debug          string
	DistanceFinder MinDistanceFinder
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
	m.cfg.height = max(int(h), minSize) + int(m.cfg.Scale*100)
	m.cfg.width = max(int(w), minSize) + int(m.cfg.Scale*100)

	tex := NewGlyph(m.cfg.width, m.cfg.height)

	for y := range m.cfg.height {
		for x := range m.cfg.width {

			xi, yi := metrics.ToFloat(x, y)
			flippedY := m.cfg.height - 1 - y

			r := m.getChannel(contours, RED, xi, yi)
			g := m.getChannel(contours, GREEN, xi, yi)
			b := m.getChannel(contours, BLUE, xi, yi)

			tex.Image().Set(x, flippedY, color.RGBA{r, g, b, 255})

		}

	}

	if m.cfg.Debug != "" {
		dbg := NewGlyph(512, 512)
		for _, con := range contours {
			con.Debug(dbg, metrics)
		}
		dbg.Save(fmt.Sprintf("%s/%c_debug.png", m.cfg.Debug, r))

	}
	return tex
}

func (m *Msdf) getChannel(contours []*Contour, c EdgeColor, x, y float64) uint8 {

	var A *Vector
	var B *Vector
	found := false
	minDistance := math.MaxFloat32

	for _, contour := range contours {
		for _, edge := range contour.Edges {

			curve := edge.Curve
			color := edge.Color

			if !color.Has(c) {
				continue
			}

			Q := Point{X: x, Y: y}
			t, d := m.cfg.DistanceFinder.Find(curve, Q)

			if d < minDistance {
				found = true
				minDistance = d
				P := curve.PointAt(t)

				A = P.Sub(Q)
				B = curve.TangentAt(t).asVector()
			}

		}
	}

	if !found {
		return 127
	}

	distance := sign(B.Cross(A)) * (minDistance)

	pixelSize := math.Min(float64(m.cfg.width), float64(m.cfg.height))
	distanceRange := (2.0 / pixelSize) * 50

	normalized := (distance / distanceRange) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
