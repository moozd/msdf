package msdf

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Edge struct {
	id    int
	Kind  string
	Color EdgeColor
	Curve *Curve
}

func (m *Msdf) getEdges(r rune) ([]*Edge, error) {
	var edges []*Edge

	segments, _, err := m.getVector(r)
	if err != nil {
		return nil, err
	}

	var p0 fixed.Point26_6

	idx := 0
	for _, segment := range segments {
		args := segment.Args
		switch segment.Op {
		case sfnt.SegmentOpMoveTo:
			p0 = args[0]
		case sfnt.SegmentOpLineTo:

			edges = append(edges, &Edge{
				id:    idx,
				Kind:  "L",
				Curve: NewCurve(&Line{P0: p0, P1: args[0]}),
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, &Edge{
				id:   idx,
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
			edges = append(edges, &Edge{
				id:   idx,
				Kind: "Q",
				Curve: NewCurve(&QuadraticBezier{
					P0: p0,
					P1: args[0],
					P2: args[1],
				}),
			})
			p0 = args[1]

		}
		idx += 1

	}

	return edges, nil

}

func (m *Msdf) getVector(r rune) (sfnt.Segments, fixed.Rectangle26_6, error) {

	ppem := fixed.I(12)

	var buff sfnt.Buffer
	gi, err := m.font.GlyphIndex(&buff, r)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, err
	}

	segments, err := m.font.LoadGlyph(&buff, gi, ppem, nil)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, err
	}

	bounds, _, err := m.font.GlyphBounds(&buff, gi, ppem, font.HintingNone)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, err
	}

	return segments, bounds, nil

}

func (e *Edge) String() string {
	return fmt.Sprintf("%s%02d[%s] ", e.Kind, e.id, e.Color)
}
