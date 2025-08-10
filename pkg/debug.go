package msdf

import (
	"image/color"
)

func (c1 *Curve) Debug(g *Glyph, color color.RGBA, scaler *Scaler) {
	img := g.Image()
	bounds := img.Bounds()

	for i, p := range c1.Points {
		pixelX, pixelY := scaler.g2p(p)

		if pixelX < bounds.Min.X || pixelX >= bounds.Max.X ||
			pixelY < bounds.Min.Y || pixelY >= bounds.Max.Y {
			continue
		}

		if i == len(c1.Points)-1 {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					x, y := pixelX+dx, pixelY+dy
					if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
						img.Set(x, y, color)
					}
				}
			}
			continue
		}

		img.Set(pixelX, pixelY, color)
	}
}
