package msdf

import (
	"fmt"
	"math"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Edges []Edge

type Edge interface {
	Has(c EdgeColor) bool
	Paint(c EdgeColor)
	Shape() Shape
}

type edge struct {
	kind  string
	color EdgeColor
	shape Shape
}

type TextureCoordScaler struct {
	bounds fixed.Rectangle26_6
	config *Config
}

type EdgeColor byte

const (
	RED   EdgeColor = 3
	GREEN           = 3 << 2
	BLUE            = 3 << 4
	CLEAR           = 0x00
)

func (e Edges) getSignedDistnace(c EdgeColor, p fixed.Point26_6) fixed.Int26_6 {
	dst := fixed.Int26_6(math.MaxInt32)
	win := 0
	edgeCount := 0
	for _, edge := range e {
		if !edge.Has(c) {
			continue
		}
		edgeCount++

		shape := edge.Shape()
		d := shape.GetDistance(p)

		if d < dst {
			dst = d
		}
		intersections := shape.RayHits(p)
		win += intersections
	}

	if edgeCount == 0 {
		return 0 // No edges found for this color
	}

	if win%2 == 1 {
		return -dst
	}

	return dst
}

func (m *Msdf) getEdges(r rune) (Edges, *TextureCoordScaler, error) {
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

	// Second pass: create edges
	for _, segment := range segments {
		args := segment.Args
		switch segment.Op {
		case sfnt.SegmentOpMoveTo:
			p0 = args[0]
		case sfnt.SegmentOpLineTo:

			edges = append(edges, &edge{
				kind:  "Line",
				shape: &LineShape{P0: p0, P1: args[0]},
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, &edge{
				kind: "Cubic",
				shape: &CubicBezierShape{
					P0: p0,
					P1: args[0],
					P2: args[1],
					P3: args[2],
				},
			})
			p0 = args[2]
		case sfnt.SegmentOpQuadTo:
			edges = append(edges, &edge{
				kind: "Quadratic",
				shape: &QuadraticBezierShape{
					P0: p0,
					P1: args[0],
					P2: args[1],
				},
			})
			p0 = args[1]

		}

	}

	// Add padding around glyph bounds for better visibility
	padding := fixed.Int26_6(2 * 64) // 2 units of padding
	bounds.Min.X -= padding
	bounds.Min.Y -= padding
	bounds.Max.X += padding
	bounds.Max.Y += padding

	// Bounds calculation working correctly

	scaler := &TextureCoordScaler{bounds: bounds, config: m.config}
	return edges, scaler, nil
}

func (e *TextureCoordScaler) scale(x, y int) fixed.Point26_6 {
	// Pure fixed-point arithmetic - no float conversions
	rangeX := e.bounds.Max.X - e.bounds.Min.X
	rangeY := e.bounds.Max.Y - e.bounds.Min.Y
	
	// Scale: (pixel / textureSize) * range + min
	// Using fixed-point multiplication/division
	fx := (fixed.Int26_6(x) * rangeX) / fixed.Int26_6(e.config.Advance) + e.bounds.Min.X
	fy := (fixed.Int26_6(y) * rangeY) / fixed.Int26_6(e.config.LineHeight) + e.bounds.Min.Y
	
	return fixed.Point26_6{
		X: fx,
		Y: fy,
	}
}

func (e *edge) Has(color EdgeColor) bool {
	return (e.color & color) == color
}

func (e *edge) Paint(color EdgeColor) {
	e.color = color
}

func (e *edge) Shape() Shape {
	return e.shape
}

func (e *edge) String() string {
	return fmt.Sprintf("%s: %s", e.kind, e.color)
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
