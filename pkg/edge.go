package msdf

import (
	"container/heap"
	"fmt"
	"image/color"
	"math"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Edges []Edge

type Edge struct {
	Kind  string
	Color EdgeColor
	Curve *Curve
}

type Scaler struct {
	bounds fixed.Rectangle26_6
	config *Config
}

type EdgeColor byte

const (
	RED   EdgeColor = 1 << 0 // = 1 (bit 0)
	GREEN           = 1 << 1 // = 2 (bit 1)
	BLUE            = 1 << 2 // = 4 (bit 2)
	CLEAR           = 0x00
)

func (e Edges) getContours() [][]*Edge {

	var contours [][]*Edge

	contours = append(contours, []*Edge{})

	cI := 0
	includeFirstEdge := true
	for i := 1; i < len(e); i += 1 {
		a, b := &e[i-1], &e[i]

		isConnected := a.Curve.IsConnected(b.Curve)

		if isConnected {

			if includeFirstEdge {
				contours[cI] = append(contours[cI], a)
				includeFirstEdge = false
			}
			contours[cI] = append(contours[cI], b)
			continue
		}

		contours = append(contours, []*Edge{})
		includeFirstEdge = true
		cI += 1
	}

	return contours
}

func (e Edges) setupColors() {

	colorIndex := 0
	cons := e.getContours()

	for i := range cons {

		// fmt.Printf("contour: %d -------- count: %d \n", i+1, len(cons[i]))
		shrpest := make(MaxHeap, 0)
		heap.Init(&shrpest)

		for j := 1; j < len(cons[i]); j += 1 {
			a := cons[i][j-1]
			b := cons[i][j]
			v, ok := a.Curve.GetSharpCorner(b.Curve, 50)
			if ok {
				// fmt.Printf("sharp corner found\n")
				heap.Push(&shrpest, &HeapItem{value: j, priority: v})
			}
		}

		colorIndex = 0
		m1, m2, m3 := 0, 1, 2

		if shrpest.Len() >= 3 {
			m1 = heap.Pop(&shrpest).(*HeapItem).value
			m2 = heap.Pop(&shrpest).(*HeapItem).value
			m3 = heap.Pop(&shrpest).(*HeapItem).value
		}

		for j := range cons[i] {
			cons[i][j].Color = colors[colorIndex]
			if j == m1 || j == m2 || j == m3 {
				colorIndex = (colorIndex + 1) % 3
			}
		}

	}
}

func (e Edges) getSignedDistnace(c EdgeColor, p fixed.Point26_6) float64 {
	dst := math.MaxFloat64
	winding := 0
	edgeCount := 0
	for _, edge := range e {
		if !edge.Color.Has(c) {
			continue
		}
		edgeCount++

		d := edge.Curve.FindMinDistance(p)

		if d < dst {
			dst = d
		}
		w := edge.Curve.Cast(p)
		winding += w
	}

	fmt.Printf("d: %f\n", dst)
	if edgeCount == 0 {
		return 0
	}

	if winding%2 == 1 {
		return -dst
	}

	return dst
}

func (m *Msdf) getEdges(r rune) (Edges, *Scaler, error) {
	var edges []Edge

	ppem := fixed.I(12)

	var buff sfnt.Buffer
	gi, err := m.font.GlyphIndex(&buff, r)
	if err != nil {
		return nil, nil, err
	}

	segments, err := m.font.LoadGlyph(&buff, gi, ppem, nil)
	if err != nil {
		return nil, nil, err
	}

	var p0 fixed.Point26_6

	// Calculate actual glyph bounds from segments
	bounds := fixed.Rectangle26_6{
		Min: fixed.Point26_6{X: fixed.Int26_6(1 << 20), Y: fixed.Int26_6(1 << 20)},
		Max: fixed.Point26_6{X: fixed.Int26_6(-1 << 20), Y: fixed.Int26_6(-1 << 20)},
	}

	// First pass: calculate bounds
	for _, segment := range segments {
		for _, arg := range segment.Args {
			if arg.X < bounds.Min.X {
				bounds.Min.X = arg.X
			}
			if arg.Y < bounds.Min.Y {
				bounds.Min.Y = arg.Y
			}
			if arg.X > bounds.Max.X {
				bounds.Max.X = arg.X
			}
			if arg.Y > bounds.Max.Y {
				bounds.Max.Y = arg.Y
			}
		}
	}

	fmt.Printf("Glyph Coords: Y(%f,%f) X(%f,%f)\n",
		unpack_i26_6(bounds.Min.Y),
		unpack_i26_6(bounds.Max.Y),
		unpack_i26_6(bounds.Min.X),
		unpack_i26_6(bounds.Max.X),
	)

	// Second pass: create edges
	for _, segment := range segments {
		args := segment.Args
		switch segment.Op {
		case sfnt.SegmentOpMoveTo:
			p0 = args[0]
		case sfnt.SegmentOpLineTo:

			edges = append(edges, Edge{
				Kind:  "L",
				Curve: NewCurve(&Line{P0: p0, P1: args[0]}),
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, Edge{
				Kind: "C",
				Curve: NewCurve(&CubicBezier{
					P0: p0,
					P1: args[0],
					P2: args[1],
					P3: args[2],
				}),
			})
			p0 = args[2]
		case sfnt.SegmentOpQuadTo:
			edges = append(edges, Edge{
				Kind: "Q",
				Curve: NewCurve(&QuadraticBezier{
					P0: p0,
					P1: args[0],
					P2: args[1],
				}),
			})
			p0 = args[1]

		}

	}

	padding := pack_i26_6(0.5)
	bounds.Min.X -= padding
	bounds.Min.Y -= padding
	bounds.Max.X += padding
	bounds.Max.Y += padding

	scaler := &Scaler{bounds: bounds, config: m.cfg}
	return edges, scaler, nil

}

func (e *Scaler) p2g(x, y int) fixed.Point26_6 {
	rangeX := e.bounds.Max.X - e.bounds.Min.X
	rangeY := e.bounds.Max.Y - e.bounds.Min.Y

	fx := (fixed.I(x)*rangeX)/fixed.I(e.config.Advance) + e.bounds.Min.X
	fy := e.bounds.Min.Y - (fixed.I(y)*rangeY)/fixed.I(e.config.LineHeight)

	return fixed.Point26_6{
		X: fx,
		Y: fy,
	}
}

func (e *Scaler) g2p(p fixed.Point26_6) (int, int) {
	rangeX := e.bounds.Max.X - e.bounds.Min.X
	rangeY := e.bounds.Max.Y - e.bounds.Min.Y

	// Convert back from glyph coords to pixel coords
	pixelX := ((p.X - e.bounds.Min.X) * fixed.I(e.config.Advance)) / rangeX
	pixelY := ((p.Y - e.bounds.Min.Y) * fixed.I(e.config.LineHeight)) / rangeY

	return pixelX.Round(), pixelY.Round()
}

func (e EdgeColor) RGB() color.RGBA {

	var r, g, b uint8

	if (e & RED) == RED {
		r = 255
	}

	if (e & GREEN) == GREEN {
		g = 255
	}

	if (e & BLUE) == BLUE {
		b = 255
	}
	return color.RGBA{r, g, b, 255}
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Color)
}

func (e EdgeColor) Has(color EdgeColor) bool {
	return (e & color) == color
}

func (e EdgeColor) String() string {
	str := ""
	isRed := e&RED == RED
	isGreen := e&GREEN == GREEN
	isBlue := e&BLUE == BLUE

	if isRed {
		str = fmt.Sprintf("%s RED", str)
	}

	if isGreen {
		str = fmt.Sprintf("%s GREEN", str)
	}

	if isBlue {
		str = fmt.Sprintf("%s BLUE", str)
	}

	if e == CLEAR {
		return "CLEAR"
	}
	return str
}
