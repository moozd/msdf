package msdf

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type Glyph struct {
	img *image.RGBA
	cfg *Config
}

func NewGlyph(cfg *Config) *Glyph {
	o := &Glyph{
		cfg: cfg,
	}

	o.img = image.NewRGBA(image.Rect(0, 0, o.cfg.Advance, o.cfg.LineHeight))
	bg := &image.Uniform{color.RGBA{0, 0, 0, 255}}
	draw.Draw(o.img, o.img.Bounds(), bg, image.Point{}, draw.Src)

	return o
}

func (o *Glyph) Save(s string) {
	file, _ := os.Create(s)
	defer file.Close()
	png.Encode(file, o.img)
}

func (o *Glyph) Image() *image.RGBA {
	return o.img
}
