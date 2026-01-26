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
	cache    map[rune]*Glyph
}

type Config struct {
	Seed             uint
	Size             float64
	DistanceField    float64
	Scale            float64
	DebugArtifactDir string
	DistanceFinder   MinDistanceFinder
	EdgeColorizer    EdgeColorizer
	height, width    int
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

	if cfg.EdgeColorizer == nil {
		cfg.EdgeColorizer = &SimpleEdgeColorizer{}
	}
	if cfg.DistanceFinder == nil {
		cfg.DistanceFinder = &BruteForceMinDistanceFinder{}
	}

	msdf := &Msdf{
		cfg:      cfg,
		font:     fnt,
		cache:    make(map[rune]*Glyph),
		metadata: &Metadata{},
		palette:  newEdgeColorPalette(&cfg.Seed),
	}

	return msdf, nil
}

func (m *Msdf) CreateAtlas(size int, charset string) (*Canvas, *Metadata, error) {
	atlas := newCanvas(size, size, color.RGBA{0, 0, 0, 255})
	m.metadata.Atlas.Type = "msdf"
	m.metadata.Atlas.YOrigin = "bottom"
	m.metadata.Atlas.Size = m.cfg.Size
	m.metadata.Atlas.Width = size
	m.metadata.Atlas.Height = size
	m.metadata.Glyphs = []GlyphOptions{}

	x := 0
	y := 0
	var glyps []*Glyph
	done := make(map[rune]bool)

	for _, c := range charset {
		if done[c] == true {
			continue
		}
		done[c] = true
		glyph, err := m.Get(c)
		if err != nil {
			return nil, nil, err
		}
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

		m.metadata.Glyphs = append(m.metadata.Glyphs, *g.Options)
	}

	return atlas, m.metadata, nil
}

func (m *Msdf) Get(r rune) (*Glyph, error) {
	glyph, ok := m.cache[r]
	if ok {
		return glyph, nil
	}
	metrics, err := m.getMetrics(r)
	if err != nil {
		return nil, err
	}

	contours, err := m.getContours(r)
	if err != nil {
		return nil, err
	}

	bounds := metrics.GetBounds()
	m.cfg.width = bounds.Dx()
	m.cfg.height = bounds.Dy()

	planeBounds := metrics.GetPlaneBounds()
	options := GlyphOptions{
		Unicode: int(r),
		Advance: metrics.GetAdvance(),
		PlaneBounds: Coords{
			Left:   unpackI26_6(planeBounds.Min.X),
			Right:  unpackI26_6(planeBounds.Max.X),
			Top:    unpackI26_6(planeBounds.Min.Y),
			Bottom: unpackI26_6(planeBounds.Max.Y),
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

	if m.cfg.DebugArtifactDir != "" {
		dbg := newCanvas(512, 512, color.RGBA{0, 0, 0, 255})
		for _, con := range contours {
			con.Debug(dbg, metrics)
		}
		dbg.Save(fmt.Sprintf("%s/%c_debug.png", m.cfg.DebugArtifactDir, r))
		m.render(r, canvas)
	}

	m.cache[r] = newGlyph(canvas, &options)
	return m.cache[r], nil
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
