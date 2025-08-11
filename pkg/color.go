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

	// each  edge will get a pair of colors
	// adjacent ones will at least have one color in common

	var current EdgeColor

	for _, contour := range contours {
		edges := contour.edges
		if len(edges) == 1 {
			current = WHITE
		} else {
			current = RED | BLUE
		}

		for _, edge := range edges {
			edge.Color = current
			if current.Has(RED | GREEN) {
				current = GREEN | BLUE
			} else {
				current = RED | GREEN
			}

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
