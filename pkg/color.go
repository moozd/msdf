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
		fmt.Println()
		fmt.Printf("Contour: %d\n", k+1)
		edges := contour.edges
		n := len(edges)

		for i := range n {
			nextIdx := (i + 1) % n
			isSharp, deg := edges[i].Curve.IsCorner(edges[nextIdx].Curve, contour.winding, 136)
			fmt.Printf("%v->%v: isSharp: %-8t angle: %-8.2f deg\n", edges[i], edges[nextIdx], isSharp, deg)
		}

		for i, edge := range edges {
			edge.Color = []EdgeColor{RED | GREEN, GREEN | BLUE, BLUE | RED}[i%3]
		}

		for _, edge := range edges {
			fmt.Printf("%v\n", edge)
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
