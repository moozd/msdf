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
	Seed          uint
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

func getDistance(c Curve, Q Point) (float64, float64) {

	queue := []float64{0.0, 0.1}

	for len(queue) > 2 {
		ts := queue[0:1][0]
		te := queue[1:2][0]
		queue = queue[2:]
		cs := c.CurvatureAt(ts).asVector().Distance()
		ce := c.CurvatureAt(te).asVector().Distance()

		if math.Abs(ce-cs) < 1e-8 {
			break
		}
		queue = append(queue, ts, te/2.0, te/2.0, te)

	}

	return 0.0, 0.0
}

func getChannel(cfg *Config, contours []*Contour, c EdgeColor, x, y float64) uint8 {

	var A *Vector
	var B *Vector
	found := false
	newtonRaphsonT := 0.0
	newtonRaphsonMinDistance := math.MaxFloat32
	broutforceT := 0.0
	brouteForceMinDistance := math.MaxFloat32
	for _, con := range contours {
		for _, edge := range con.Edges {
			curve := edge.Curve
			if !edge.Color.Has(c) {
				continue
			}

			for t := 0.0; t <= 1; t += 0.1 {

				p := curve.PointAt(t)
				a := vec().fromXY(p.X, p.Y, x, y)

				d := a.Distance()

				if d < brouteForceMinDistance {
					brouteForceMinDistance = d
					broutforceT = t
				}
			}

			// FIXME: fix  Newton Raphson
			d, t := getDistance(curve, Point{X: x, Y: y})
			if d < newtonRaphsonMinDistance {
				found = true
				newtonRaphsonT = t
				newtonRaphsonMinDistance = d
				p := curve.PointAt(t)
				A = vec().fromXY(p.X, p.Y, x, y)
				B = vec().fromP(curve.TangentAt(t))
			}
		}
	}

	// Debug: always print comparison for first few points to verify fix
	if x < 2.4 && y > -0.5 {
		fmt.Printf("X: %f Y: %f\n", x, y)
		fmt.Printf("BF: T=%f D=%f\n", broutforceT, brouteForceMinDistance)
		fmt.Printf("NR: T=%f D=%f\n", newtonRaphsonT, newtonRaphsonMinDistance)
		fmt.Printf("NR T in range [0,1]: %t\n", 0 <= newtonRaphsonT && newtonRaphsonT <= 1)
		fmt.Println()
	}

	if !found {
		return 127
	}

	distance := sign(B.Cross(A)) * (newtonRaphsonMinDistance)

	pixelSize := math.Min(float64(cfg.width), float64(cfg.height))
	distanceRange := (2.0 / pixelSize) * 50

	normalized := (distance / distanceRange) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
