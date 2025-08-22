package msdf

import (
	"image/color"
	"log"
)

type EdgeColorizer interface {
	Colorize(contours []*Contour, palette *EdgeColorPalette)
}

type SimpleEdgeColorizer struct{}

func (sc SimpleEdgeColorizer) verbose(c string, contour *Contour) {
	log.Printf("Colorizer: case=%s char=%c\n", c, contour.Symbol)
	for _, e := range contour.Edges {
		log.Printf("Colorizer: edge=%v %v\n", e, e.Curve)
	}
}

func (sc SimpleEdgeColorizer) Colorize(contours []*Contour, palette *EdgeColorPalette) {

	thershold := 126 * RAD
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

		if len(corners) == 0 { // smooth edge
			sc.verbose("smooth:before", contour)

			palette.Shuffle(&color)
			for _, edge := range edges {
				edge.Color = color
			}

			sc.verbose("smooth:after", contour)

		} else if len(corners) == 1 { // teardrop case
			sc.verbose("teardrop:before", contour)
			var colors [3]EdgeColor
			palette.Shuffle(&color)
			colors[0] = color
			colors[1] = WHITE
			palette.Shuffle(&color)
			colors[2] = color

			E := len(contour.Edges)
			corner := corners[0]
			if E >= 3 {
				for i := range contour.Edges {
					contour.Edges[(corner+i)%E].Color = colors[1+symmetricalMagic(i, E)]
				}
			} else if E >= 1 {
				contour.Ensure3Edges()

				if E >= 2 {
					contour.Edges[0].Color = colors[0]
					contour.Edges[1].Color = colors[0]
					contour.Edges[2].Color = colors[1]
					contour.Edges[3].Color = colors[1]
					contour.Edges[4].Color = colors[2]
					contour.Edges[5].Color = colors[2]
				} else {

					contour.Edges[0].Color = colors[0]
					contour.Edges[1].Color = colors[1]
					contour.Edges[2].Color = colors[2]
				}
			}

			sc.verbose("teardrop:after", contour)

		} else { // multiple corners
			sc.verbose("multiple:before", contour)
			cornerCount := len(corners)
			spline := 0
			start := corners[0]
			palette.Shuffle(&color)
			initialColor := color
			for i := range contour.Edges {
				index := (start + i) % len(contour.Edges)
				if spline+1 < cornerCount && corners[spline+1] == index {
					spline += 1
					banned := EdgeColor(0)
					if spline == cornerCount-1 {
						banned = initialColor
					}
					palette.ShuffleEx(&color, EdgeColor(banned))
				}
				contour.Edges[index].Color = color
			}

			sc.verbose("multiple:after", contour)
		}
	}

}

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

type EdgeColorPalette struct {
	seed   *uint
	colors []EdgeColor
}

func newEdgeColorPalette(seed *uint) *EdgeColorPalette {
	return &EdgeColorPalette{
		seed:   seed,
		colors: []EdgeColor{CYAN, MAGENTA, YELLOW},
	}
}

func (cp *EdgeColorPalette) seedExtract2() int {

	v := int(*cp.seed) & 1
	*cp.seed = *cp.seed >> 1
	return v
}

func (cp *EdgeColorPalette) seedExtract3() int {
	v := int(*cp.seed % 3)
	*cp.seed /= 3
	return v
}

func (cp *EdgeColorPalette) Init() EdgeColor {
	c := cp.colors[cp.seedExtract3()]
	log.Printf("EdgeColorPalette: Initialize color ->  (%v)", c)
	return c
}

func (cp *EdgeColorPalette) Shuffle(color *EdgeColor) {
	old := *color

	shifted := *color << (1 + cp.seedExtract2())
	*color = EdgeColor((shifted | shifted>>3) & WHITE)

	log.Printf("EdgeColorPalette: Shuffle (%v) -> (%v)", old, *color)
}

func (cp *EdgeColorPalette) ShuffleEx(color *EdgeColor, banned EdgeColor) {
	old := *color

	combined := EdgeColor(*color & banned)
	if combined == RED || combined == GREEN || combined == BLUE {
		*color = EdgeColor(combined ^ WHITE)
	} else {
		cp.Shuffle(color)
	}

	log.Printf("EdgeColorPalette: ShuffleEx (%v) -> (%v)", old, *color)
}

// the original function is called symmetricalTrichotomy in Chlumsky's implementation
// I have no idea how and why it works, but it just selects a random index using pos and n
func symmetricalMagic(pos, n int) int {
	return int(3.0+2.875*float64(pos)/(float64(n)-1.0)-1.4375+.5) - 3
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

	if e&WHITE == WHITE {
		return "W"
	}

	if e&MAGENTA == MAGENTA {
		return "M"
	}

	if e&YELLOW == YELLOW {
		return "Y"
	}

	if e&CYAN == CYAN {
		return "C"
	}

	if e&RED == RED {
		return "R"
	}

	if e&GREEN == GREEN {
		return "G"
	}

	if e&BLUE == BLUE {
		return "B"
	}

	if e == CLEAR {
		return "-"
	}

	return "?"
}
