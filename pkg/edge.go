package msdf

import (
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

func (e Edges) setupColors() {
	for i := range e {
		e[i].Color = colors[i%3]
	}

	// corners := make(map[*Edge][]*Edge)
	//
	// for i := range e {
	// 	corners[&e[i]] = make([]*Edge, 0)
	// 	for j := range e {
	// 		curveA := e[i].Curve
	// 		curveB := e[j].Curve
	//
	// 		isCorner := curveA.IsAttached(curveB)
	// 		if isCorner {
	// 			corners[&e[i]] = append(corners[&e[i]], &e[j])
	// 		}
	// 	}
	// }
	//
	// for i := range len(e) {
	// 	cl := colors[i%len(colors)]
	// 	// for j := range    {
	// 	// 	corners[i][j].Color = cl
	// 	// }
	// }
	//
	// for i := range e {
	// 	fmt.Println(corners[i])
	// }
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
				Kind:  "Line",
				Curve: NewCurve(&Line{P0: p0, P1: args[0]}),
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, Edge{
				Kind: "Cubic",
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
				Kind: "Quadratic",
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

	scaler := &Scaler{bounds: bounds, config: m.config}
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
