package main

import (
	"image/png"
	"os"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 256,
		Advance:    256,
	}
	tex, _ := msdf.New("/home/mo/.local/share/fonts/Monaspace/MonaspiceArNerdFont-Regular.otf", cfg)

	glyph := tex.Get('D')

	file, _ := os.Create("output.png")
	defer file.Close()
	png.Encode(file, glyph)

}
