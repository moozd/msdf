package msdf

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func (c *Contour) Debug(g *Glyph, m *Metrics) {
	count := 0
	for _, edge := range c.Edges {

		points := edge.Curve.GetLowResPoints()
		c := edge.Color.RGB()
		img := g.Image()
		bounds := img.Bounds()

		for i, p := range points {
			px, py := m.Scale(p, bounds, 40)

			if px < bounds.Min.X || px >= bounds.Max.X ||
				py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}

			// Draw label at center of curve
			if i == len(points)/2 {
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(color.White),
					Face: basicfont.Face7x13,
					Dot:  fixed.Point26_6{X: fixed.I(px), Y: fixed.I(py + 15)},
				}

				d.DrawString(fmt.Sprintf("%s%d", edge.Kind, edge.id))
				count += 1
			}

			if i == len(points)-1 {
				for dx := -2; dx <= 2; dx++ {
					for dy := -2; dy <= 2; dy++ {
						x, y := px+dx, py+dy
						if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
							img.Set(x, y, c)
						}
					}
				}
				continue
			}

			img.Set(px, py, c)
		}
	}
}

func (m *Msdf) Debug(path string, tex *Glyph) {
	out := NewGlyph(512, 512)

	texBounds := tex.Image().Bounds()

	for y := range 512 {
		for x := range 512 {
			srcX := int(float64(x) * float64(texBounds.Dx()) / 512.0)
			srcY := int(float64(y) * float64(texBounds.Dy()) / 512.0)

			if srcX < texBounds.Max.X && srcY < texBounds.Max.Y {
				r, g, b, _ := tex.Image().At(srcX, srcY).RGBA()

				channels := []uint32{r, g, b}
				if channels[0] > channels[1] {
					channels[0], channels[1] = channels[1], channels[0]
				}
				if channels[1] > channels[2] {
					channels[1], channels[2] = channels[2], channels[1]
				}
				if channels[0] > channels[1] {
					channels[0], channels[1] = channels[1], channels[0]
				}

				median := uint8(channels[1] >> 8)
				out.Image().Set(x, y, color.RGBA{median, median, median, 255})
			}
		}
	}

	out.Save(path)

}
