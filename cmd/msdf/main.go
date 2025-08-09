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
	tex, _ := msdf.New("/home/mo/.local/share/fonts/Monaspace/MonaspiceArNerdFont-Regular.otf", cfg)

	glyph := tex.Get('A')

	file, _ := os.Create("output.png")
	defer file.Close()
	png.Encode(file, glyph)

}
