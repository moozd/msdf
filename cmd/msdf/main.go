package main

import (
	"image/png"
	"os"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 128,
		Advance:    128,
	}
	tex, _ := msdf.New("/home/mo/.local/share/fonts/Monaspace/MonaspiceArNerdFont-Regular.otf", cfg)

	glyph := tex.Get('R')

	file, _ := os.Create("output.png")
	defer file.Close()
	png.Encode(file, glyph)

}
