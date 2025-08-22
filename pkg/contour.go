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
	Symbol  rune
	Edges   []*Edge
	Winding ClockDirection
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

		cons = append(cons, newContour(ce, r))

	}

	m.cfg.EdgeColorizer.Colorize(cons, m.palette)

	return cons, nil
}

func newContour(edges []*Edge, symbol rune) *Contour {

	signedAreas := 0.0
	for _, edge := range edges {
		signedAreas += edge.Curve.GetSignedArea()
	}

	w := CCW
	if signedAreas > 0 {
		w = CW
	}

	return &Contour{
		Symbol:  symbol,
		Edges:   edges,
		Winding: w,
	}
}

func (c *Contour) Ensure3Edges() {
	E := len(c.Edges)
	if E >= 3 || E < 1 {
		return
	}

	if E == 2 {
		edge := c.Edges[1]
		c1, c2 := edge.Curve.Split()
		c.Edges = []*Edge{}
		c.Edges = append(c.Edges, edge.make(1, c1), edge.make(2, c2))
	}

	if E == 1 {
		edge := c.Edges[0]
		c1, c2 := edge.Curve.Split()
		c.Edges = []*Edge{}
		c.Edges = append(c.Edges, edge.make(1, c1))
		edge = edge.make(2, c2)
		c1, c2 = edge.Curve.Split()
		c.Edges = append(c.Edges, edge.make(1, c1), edge.make(2, c2))
	}

}

func (c Contour) String() string {
	return fmt.Sprintf("D: %d , E: %v", c.Winding, c.Edges)
}

func (c ClockDirection) String() string {
	if c == CW {
		return "CW"
	}
	return "CCW"
}
