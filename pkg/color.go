package msdf

import (
	"fmt"
	"image/color"
	"strings"
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

	for k, contour := range contours {

		fmt.Printf("Contour: %d\n", k+1)
		edges := contour.edges
		for i, edge := range edges {
			nextEdge := edges[(i+1)%len(edges)]

			// Check if this edge is part of a sharp corner
			isSharpCorner := edge.Curve.IsCorner(nextEdge.Curve, contour.winding, 136)

			if isSharpCorner {
				fmt.Printf("%v , %v  %v\n", edge, nextEdge, contour.winding)
			}

			// if isSharpCorner {
			// 	edge.Color = []EdgeColor{RED | BLUE, GREEN | BLUE, RED | GREEN}[i%3]
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
	isRed := e&RED == RED
	isGreen := e&GREEN == GREEN
	isBlue := e&BLUE == BLUE

	colors := []string{"-", "-", "-"}

	if isRed {
		colors[0] = "R"
	}

	if isGreen {
		colors[1] = "G"
	}

	if isBlue {
		colors[2] = "B"
	}

	return strings.Join(colors, " ")
}
