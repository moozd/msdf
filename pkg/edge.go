package msdf

import (
	"fmt"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Edge struct {
	Kind  string
	Color EdgeColor
	Curve *Curve
}

func (m *Msdf) getEdges(r rune) ([]*Edge, error) {
	var edges []*Edge

	segments, err := m.getSegments(r)
	if err != nil {
		return nil, err
	}

	var p0 fixed.Point26_6

	for _, segment := range segments {
		args := segment.Args
		switch segment.Op {
		case sfnt.SegmentOpMoveTo:
			p0 = args[0]
		case sfnt.SegmentOpLineTo:

			edges = append(edges, &Edge{
				Kind:  "L",
				Curve: NewCurve(&Line{P0: p0, P1: args[0]}),
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, &Edge{
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

	return edges, nil

}

func (m *Msdf) getSegments(r rune) (sfnt.Segments, error) {

	ppem := fixed.I(12)

	var buff sfnt.Buffer
	gi, err := m.font.GlyphIndex(&buff, r)
	if err != nil {
		return nil, err
	}

	segments, err := m.font.LoadGlyph(&buff, gi, ppem, nil)
	if err != nil {
		return nil, err
	}

	return segments, nil

}

func (e *Edge) String() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Color)
}
