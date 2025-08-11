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
	CLEAR           = 0x00
)

func colorize(cons []*Contour) {

	// each  edge will get a pair of colors
	// adjacent ones will atleast have one color in common
	colors := []EdgeColor{RED, GREEN, BLUE}
	ci := 0

	for i := range cons {
		edges := cons[i].edges
		for j := 1; j < len(edges); j += 1 {
			a := edges[j-1]
			b := edges[j]
			ci = ci + 1
			seed1 := colors[ci%3]
			seed2 := colors[(ci+1)%3]
			seed3 := colors[(ci+2)%3]

			if a.Curve.IsConnected(b.Curve) {
				a.Color = seed1 | seed2
				b.Color = seed1 | seed3
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
