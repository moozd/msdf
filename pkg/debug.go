package msdf

import (
	"image/color"
)

func (c1 *Curve) Debug(g *Glyph, color color.RGBA, scaler *Metrics) {
	img := g.Image()
	bounds := img.Bounds()

	for i, p := range c1.Points {
		px, py := scaler.ToPixel(p)

		if px < bounds.Min.X || px >= bounds.Max.X ||
			py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}

		if i == len(c1.Points)-1 {
			for dx := -1; dx <= 1; dx++ {
				for dy := -4; dy <= 4; dy++ {
					x, y := px+dx, py+dy
					if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
						img.Set(x, y, color)
					}
				}
			}
			continue
		}

		img.Set(px, py, color)
	}
}
