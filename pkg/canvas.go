package msdf

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type Canvas struct {
	buffer *image.RGBA
}

func newCanvas(width, height int, bg color.RGBA) *Canvas {
	c := &Canvas{
		buffer: image.NewRGBA(image.Rect(0, 0, width, height)),
	}

	bgUniform := &image.Uniform{bg}
	draw.Draw(c.buffer, c.buffer.Bounds(), bgUniform, image.Point{}, draw.Src)

	return c
}

func (c *Canvas) Set(x, y int, r, g, b uint8) {
	c.buffer.Set(x, y, color.RGBA{r, g, b, 255})
}

func (c *Canvas) Put(x, y int, g *Canvas) {
	srcBounds := c.buffer.Bounds()
	flippedY := y
	dstRect := image.Rect(x, flippedY, x+srcBounds.Dx(), y+srcBounds.Dy())
	draw.Draw(c.buffer, dstRect, g.buffer, srcBounds.Min, draw.Over)
}

func (c *Canvas) Save(s string) {
	file, _ := os.Create(s)
	defer file.Close()
	png.Encode(file, c.buffer)
}

func (c *Canvas) Image() *image.RGBA {
	return c.buffer
}
