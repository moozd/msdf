package msdf

import (
	"fmt"
	"image/color"
)

type EdgeColor byte

const (
	RED   EdgeColor = 1 << 0 // = 1 (bit 0)
	GREEN           = 1 << 1 // = 2 (bit 1)
	BLUE            = 1 << 2 // = 4 (bit 2)
	WHITE           = RED | GREEN | BLUE
	CLEAR           = 0x00
)

func colorize(contours []*Contour) {

	for _, contour := range contours {

		edges := contour.edges
		for i, edge := range edges {
			// nextEdge := edges[(i+1)%len(edges)]
			//
			// // Check if this edge is part of a sharp corner
			// isSharpCorner := edge.Curve.IsCorner(nextEdge.Curve, 135)
			//
			// if isSharpCorner {
			// 	edge.Color = []EdgeColor{RED, GREEN, BLUE}[i%3]
			// } else {
			edge.Color = []EdgeColor{RED | GREEN, RED | BLUE, GREEN | BLUE}[i%3]
			// }
		}
	}

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
