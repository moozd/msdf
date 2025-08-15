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
	Scale         float64
	Debug         string
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

			r := getChannel(m.cfg, contours, RED, xi, yi)
			g := getChannel(m.cfg, contours, GREEN, xi, yi)
			b := getChannel(m.cfg, contours, BLUE, xi, yi)

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

func getChannel(cfg *Config, contours []*Contour, c EdgeColor, x, y float64) uint8 {

	var A *Vector
	var B *Vector
	found := false
	minDist := math.MaxFloat32
	for _, con := range contours {
		for _, edge := range con.edges {
			curve := edge.Curve
			if !edge.Color.Has(c) {
				continue
			}

			for t := 0.0; t <= 1; t += 0.01 {

				p := curve.PointAt(t)
				a := vec(p.X, p.Y, x, y)

				if a.Distance() < minDist {
					found = true
					minDist = a.Distance()

					A = a
					s := curve.TangentAt(t)
					B = vec(0, 0, s.X, s.Y)
				}
			}

		}
	}

	if !found {
		return 127
	}

	distance := sign(B.Cross(A)) * (minDist)

	pixelSize := math.Min(float64(cfg.width), float64(cfg.height))
	distanceRange := (2.0 / pixelSize) * 120

	normalized := (distance / distanceRange) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
