package msdf

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"sort"

	"golang.org/x/image/font/sfnt"
)

type Msdf struct {
	font     *sfnt.Font
	cfg      *Config
	metadata *Metadata
	palette  *EdgeColorPalette
}

type Config struct {
	Seed           uint
	Size           float64
	DistanceField  float64
	Scale          float64
	Debug          string
	DistanceFinder MinDistanceFinder
	EdgeColorizer      EdgeColorizer
	height, width  int
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
		cfg:      cfg,
		font:     fnt,
		metadata: &Metadata{},
		palette:  newEdgeColorPalette(&cfg.Seed),
	}

	return msdf, nil
}

func (m *Msdf) CreateAtlas(size int) (*Canvas, Metadata) {
	atlas := newCanvas(size, size, color.RGBA{0, 0, 0, 255})
	m.metadata.Altas.Type = "msdf"
	m.metadata.Altas.YOrigin = "bottom"
	m.metadata.Altas.Size = m.cfg.Size
	m.metadata.Altas.Width = size
	m.metadata.Altas.Height = size
	m.metadata.Glyphs = []GlyphOptions{}

	x := 0
	y := 0
	var glyps []*Glyph

	for i := 32; i < 128; i++ {
		c := rune(i)
		glyph := m.Get(c)
		glyps = append(glyps, glyph)
	}

	sort.Slice(glyps, func(i, j int) bool {
		ay := glyps[i].Canvas.Image().Bounds().Dy()
		by := glyps[j].Canvas.Image().Bounds().Dy()
		return (ay > by)
	})

	rowHeight := 0
	for _, g := range glyps {
		canvas := g.Canvas
		img := canvas.Image()

		currentHeight := img.Bounds().Dy()
		if currentHeight > rowHeight {
			rowHeight = currentHeight
		}

		if x+img.Bounds().Dx() > size {
			x = 0
			y += rowHeight
			rowHeight = 0
		}

		atlas.Put(x, y, canvas)
		g.setUV(Coords{
			Left:   float64(x + img.Bounds().Min.X),
			Right:  float64(x + img.Bounds().Max.X),
			Top:    float64(y + img.Bounds().Min.Y),
			Bottom: float64(y + img.Bounds().Max.Y),
		})

		x += img.Bounds().Dx()

	}

	return atlas, *m.metadata
}

func (m *Msdf) Get(r rune) *Glyph {
	metrics, _ := m.getMetrics(r)
	contours, _ := m.getContours(r)
	bounds := metrics.GetBounds()
	m.cfg.width = bounds.Dx()
	m.cfg.height = bounds.Dy()

	planeBounds := metrics.GetPlaneBounds()
	options := GlyphOptions{
		Unicode: int(r),
		Advance: metrics.GetAdvance(),
		PlaneBounds: Coords{
			Left:   unpack_i26_6(planeBounds.Min.X),
			Right:  unpack_i26_6(planeBounds.Max.X),
			Top:    unpack_i26_6(planeBounds.Min.Y),
			Bottom: unpack_i26_6(planeBounds.Max.Y),
		},
	}

	canvas := newCanvas(m.cfg.width, m.cfg.height, color.RGBA{0, 0, 0, 255})

	for y := range m.cfg.height {
		for x := range m.cfg.width {

			xi, yi := float64(bounds.Min.X+x), float64(bounds.Min.Y+y)

			r := m.getChannel(contours, RED, xi, yi)
			g := m.getChannel(contours, GREEN, xi, yi)
			b := m.getChannel(contours, BLUE, xi, yi)

			canvas.Set(x, y, r, g, b)

		}

	}

	if m.cfg.Debug != "" {
		dbg := newCanvas(512, 512, color.RGBA{0, 0, 0, 255})
		for _, con := range contours {
			con.Debug(dbg, metrics)
		}
		dbg.Save(fmt.Sprintf("%s/%c_debug.png", m.cfg.Debug, r))

	}
	return newGlyph(canvas, &options)
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

	normalized := (distance / m.cfg.DistanceField) + 0.5
	clamped := clamp(normalized, 0, 1)

	return uint8(clamped * 255)
}
