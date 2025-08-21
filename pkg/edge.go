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
	Curve Curve
}

func (m *Msdf) getEdges(r rune) ([]*Edge, error) {
	var edges []*Edge

	segments, _, _, err := m.getVector(r)
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
				Curve: newLine(p0, args[0]),
			})
			p0 = args[0]
		case sfnt.SegmentOpCubeTo:
			edges = append(edges, &Edge{
				id:    idx,
				Kind:  "C",
				Curve: newCubicBezier(p0, args[0], args[1], args[2]),
			})
			p0 = args[2]
		case sfnt.SegmentOpQuadTo:
			edges = append(edges, &Edge{
				id:    idx,
				Kind:  "Q",
				Curve: newQuadraticBezier(p0, args[0], args[1]),
			})
			p0 = args[1]

		}
		idx += 1

	}

	return edges, nil

}

func (m *Msdf) getVector(r rune) (sfnt.Segments, fixed.Rectangle26_6, fixed.Int26_6, error) {

	ppem := pack_i26_6(m.cfg.Size)

	var buff sfnt.Buffer
	gi, err := m.font.GlyphIndex(&buff, r)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, 0, err
	}

	segments, err := m.font.LoadGlyph(&buff, gi, ppem, nil)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, 0, err
	}

	bounds, adv, err := m.font.GlyphBounds(&buff, gi, ppem, font.HintingNone)
	if err != nil {
		return nil, fixed.Rectangle26_6{}, 0, err
	}

	return segments, bounds, adv, nil

}
func (e *Edge) ID() string {
	return fmt.Sprintf("%s%d", e.Kind, e.id)
}

func (e *Edge) String() string {
	return fmt.Sprintf("%s%02d[%s] ", e.Kind, e.id, e.Color)
}
