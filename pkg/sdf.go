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

	candidates := [][]float64{
		{0, vec().fromAB(Q, c.PointAt(0)).Distance()},
		{1, vec().fromAB(Q, c.PointAt(1)).Distance()},
	}
	tStarts := []float64{0.1, 0.3, 0.5, 0.9}

	for _, ts := range tStarts {
		t := solveNewtonRaphson(c, Q, ts)
		if 0 <= t && t <= 1 {
			d := vec().fromAB(Q, c.PointAt(t)).Distance()
			candidates = append(candidates, []float64{t, d})
		} else {

			for _, s := range []float64{t - 0.1, t, t + 0.1} {
				if 0 <= s && s <= 1 {
					d := vec().fromAB(Q, c.PointAt(s)).Distance()
					candidates = append(candidates, []float64{s, d})
				}
			}
		}
	}

	m := math.MaxFloat32
	res := []float64{0, 0}
	for _, c := range candidates {
		if c[1] < m {
			m = c[1]
			res = c
		}
	}
	return res[0], res[1]

}

func solveNewtonRaphson(c Curve, q Point, ts float64) float64 {

	t := ts
	for _ = range 10 {

		Q := vec().fromP(q)
		BT := vec().fromP(c.PointAt(t))
		BPT := vec().fromP(c.TangentAt(t))
		BPPT := vec().fromP(c.CurvatureAt(t))

		ft := BT.Sub(Q).Dot(BPT)
		fpt := math.Pow(BPT.Distance(), 2) + BT.Sub(Q).Dot(BPPT)

		t0 := t
		t = t - ft/fpt

		if math.Abs(t-t0) < 1e-8 {
			break
		}
	}

	return t //vec().fromAB(q, c.PointAt(t)).Distance()
}

func getChannel(cfg *Config, contours []*Contour, c EdgeColor, x, y float64) uint8 {

	var A *Vector
	var B *Vector
	found := false
	minDist := math.MaxFloat32
	for _, con := range contours {
		for _, edge := range con.Edges {
			curve := edge.Curve
			if !edge.Color.Has(c) {
				continue
			}

			for t := 0.0; t <= 1; t += 0.01 {

				p := curve.PointAt(t)
				a := vec().fromXY(p.X, p.Y, x, y)

				d := A.Distance()

				if d < minDist {
					found = true
					minDist = d
					A = a
					B = vec().fromP(curve.TangentAt(t))
				}
			}

			// FIXME: fix  Newton Raphson
			// d, t := getDistance(curve, Point{X: x, Y: y})
			// if d < minDist {
			// 	found = true
			// 	minDist = d
			// 	p := curve.PointAt(t)
			// 	A = vec().fromXY(p.X, p.Y, x, y)
			// 	B = vec().fromP(curve.TangentAt(t))
			// }

		}
	}

	if !found {
		return 127
	}

	distance := sign(B.Cross(A)) * (minDist)

	pixelSize := math.Min(float64(cfg.width), float64(cfg.height))
	distanceRange := (2.0 / pixelSize) * 50

	normalized := (distance / distanceRange) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
