package msdf

import (
	"image/color"
	"strings"
)

type EdgeColorizer interface {
	Colorize(contours []*Contour, palette *EdgeColorPalette)
}

type SimpleEdgeColorizer struct{}

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
			palette.Shuffle(&color)
			for _, edge := range edges {
				edge.Color = color
			}

		} else if len(corners) == 1 { // teardrop case
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

		} else { // multiple corners
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
	return cp.colors[cp.seedExtract3()]
}

func (cp *EdgeColorPalette) Shuffle(color *EdgeColor) {
	shifted := *color << (1 + cp.seedExtract2())
	*color = EdgeColor((shifted | shifted>>3) & WHITE)
}

func (cp *EdgeColorPalette) ShuffleEx(color *EdgeColor, banned EdgeColor) {
	combined := EdgeColor(*color & banned)
	if combined == RED || combined == GREEN || combined == BLUE {
		*color = EdgeColor(combined ^ WHITE)
	} else {
		cp.Shuffle(color)
	}
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
