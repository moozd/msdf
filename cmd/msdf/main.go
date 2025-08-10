package main

import (
	"image/png"
	"os"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 512,
		Advance:    512,
	}
	tex, _ := msdf.New("/Users/mohammad.mohammadzade/Library/Fonts/FiraCodeNerdFont-Regular.ttf", cfg)

	glyph := tex.Get('*')

	file, _ := os.Create("output.png")
	defer file.Close()
	png.Encode(file, glyph)

}
