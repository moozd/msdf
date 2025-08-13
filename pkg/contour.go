package msdf

import (
	"fmt"
)

type ClockDirection int

const (
	CW  ClockDirection = 1
	CCW ClockDirection = 0
)

type Contour struct {
	winding ClockDirection
	edges   []*Edge
}

func (m *Msdf) getContours(r rune) ([]*Contour, error) {
	edges, err := m.getEdges(r)
	if err != nil {
		return nil, err
	}

	var cons []*Contour
	var bag []*Edge

	var a, b *Edge

	for i := range edges {
		isConnected := false
		a = edges[i]

		bag = append(bag, a)
		if i+1 < len(edges) {
			b = edges[i+1]
			isConnected = a.Curve.IsConnected(b.Curve)
		}

		if isConnected {
			continue
		}

		ce := make([]*Edge, len(bag))
		copy(ce, bag)
		bag = nil

		cons = append(cons, newContour(ce))

	}

	colorize(cons)

	return cons, nil
}

func newContour(edges []*Edge) *Contour {

	signedAreas := 0.0
	for _, edge := range edges {
		signedAreas += edge.Curve.GetSignedArea()
	}

	w := CCW
	if signedAreas > 0 {
		w = CW
	}

	return &Contour{
		edges:   edges,
		winding: w,
	}
}

func (c Contour) String() string {
	return fmt.Sprintf("D: %d , E: %v", c.winding, c.edges)
}

func (c ClockDirection) String() string {
	if c == CW {
		return "CW"
	}
	return "CCW"
}
