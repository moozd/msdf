package msdf

import (
	"image"
	"image/color"
)

func (c *Curve) Debug(img *image.RGBA, color color.RGBA, scaler *Scaler) {
	bounds := img.Bounds()

	for i, p := range c.Points {
		pixelX, pixelY := scaler.g2p(p)

		if pixelX < bounds.Min.X || pixelX >= bounds.Max.X ||
			pixelY < bounds.Min.Y || pixelY >= bounds.Max.Y {
			continue
		}

		if i == len(c.Points)-1 {
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
