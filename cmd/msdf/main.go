package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 512,
		Advance:    512,
		Debug:      true,
	}

	bank := []rune{'A', 'R', '@', 'C', 'B'}

	tex, _ := msdf.New("/home/mo/.local/share/fonts/Hack/HackNerdFont-Regular.ttf", cfg)
	for _, c := range bank {

		glyph := tex.Get(c)
		glyph.Save(fmt.Sprintf("assets/%c.png", c))
	}

}
