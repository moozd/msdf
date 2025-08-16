package msdf

import (
	"image/color"
	"math"
	"strings"
)

type EdgeColor byte

const (
	RED     EdgeColor = 1 << 0 // = 1 (bit 0)
	GREEN             = 1 << 1 // = 2 (bit 1)
	BLUE              = 1 << 2 // = 4 (bit 2)
	WHITE             = RED | GREEN | BLUE
	CYAN              = GREEN | BLUE
	MAGENTA           = RED | BLUE
	YELLOW            = RED | GREEN
	CLEAR             = 0x00
)

var pallete = []EdgeColor{CYAN, MAGENTA, YELLOW}

func colorize(contours []*Contour, seed uint) {
	thershold := math.Sin(3)
	var corners []int
	color := initColor(&seed)

	for _, contour := range contours {
		edges := contour.Edges

		corners = []int{}
		prevEdge := edges[0]
		for i, edge := range edges {
			isCorner := prevEdge.Curve.IsCorner(edge.Curve, thershold)
			if isCorner {
				corners = append(corners, i)
			}

			prevEdge = edge
		}

		// smooth edge
		if len(corners) == 0 {
			// TODO: smooth edge case

			// teardrop case
		} else if len(corners) == 1 {
			// TODO: add teardrp case

			// multiple corners
		} else {
			cornerCount := len(contours)
			spline := 0
			start := corners[0]
			m := len(contour.Edges)
			switchColor(&color, &seed)
			initialColor := color
			for i := range m {
				index := (start + i) % m
				if spline+1 < cornerCount && corners[spline+1] == index {
					spline += 1
					banned := EdgeColor(0)
					if spline == cornerCount-1 {
						banned = initialColor
					}
					switchColorEx(&color, &seed, EdgeColor(banned))
				}
				contour.Edges[index].Color = color
			}
		}

	}

}

func seedExtract2(seed *uint) int {
	v := int(*seed) & 1
	*seed = *seed >> 1
	return v
}

func seedExtract3(seed *uint) int {
	v := int(*seed % 3)
	*seed /= 3
	return v
}

func initColor(seed *uint) EdgeColor {
	return pallete[seedExtract3(seed)]
}

func switchColor(color *EdgeColor, seed *uint) {
	shifted := *color << (1 + seedExtract2(seed))
	*color = EdgeColor((shifted | shifted>>3) & WHITE)
}

func switchColorEx(color *EdgeColor, seed *uint, banned EdgeColor) {
	combined := EdgeColor(*color & banned)
	if combined == RED || combined == GREEN || combined == BLUE {
		*color = EdgeColor(combined ^ WHITE)
	} else {
		switchColor(color, seed)
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
