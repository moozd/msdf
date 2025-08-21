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

type Colorizer interface {
	Colorize(contours []*Contour, palette *ColorPalette)
}

type SimpleColorizer struct{}

func (sc SimpleColorizer) Colorize(contours []*Contour, palette *ColorPalette) {

	thershold := math.Sin(3)
	var corners []int
	color := palette.Init()

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
			cornerCount := len(corners)
			spline := 0
			start := corners[0]
			m := len(contour.Edges)
			palette.Shuffle(&color)
			initialColor := color
			for i := range m {
				index := (start + i) % m
				if spline+1 < cornerCount && corners[spline+1] == index {
					spline += 1
					banned := EdgeColor(0)
					if spline == cornerCount-1 {
						banned = initialColor
					}
					palette.SuffleEx(&color, EdgeColor(banned))
				}
				contour.Edges[index].Color = color
			}
		}

	}
}

type ColorPalette struct {
	seed   *uint
	colors []EdgeColor
}

func newColorPalette(seed *uint) *ColorPalette {
	return &ColorPalette{
		seed:   seed,
		colors: []EdgeColor{CYAN, MAGENTA, YELLOW},
	}
}

func (cp *ColorPalette) seedExtract2() int {

	v := int(*cp.seed) & 1
	*cp.seed = *cp.seed >> 1
	return v
}

func (cp *ColorPalette) seedExtract3() int {
	v := int(*cp.seed % 3)
	*cp.seed /= 3
	return v
}

func (cp *ColorPalette) Init() EdgeColor {
	return cp.colors[cp.seedExtract3()]
}

func (cp *ColorPalette) Shuffle(color *EdgeColor) {
	shifted := *color << (1 + cp.seedExtract2())
	*color = EdgeColor((shifted | shifted>>3) & WHITE)
}

func (cp *ColorPalette) SuffleEx(color *EdgeColor, banned EdgeColor) {
	combined := EdgeColor(*color & banned)
	if combined == RED || combined == GREEN || combined == BLUE {
		*color = EdgeColor(combined ^ WHITE)
	} else {
		cp.Shuffle(color)
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
