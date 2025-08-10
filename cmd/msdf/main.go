package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 500,
		Advance:    500,
	}

	C := 'R'

	tex, _ := msdf.New("/home/mo/.local/share/fonts/Hack/HackNerdFont-Regular.ttf", cfg)

	glyph := tex.Get(C)

	glyph.Save(fmt.Sprintf("%c.png", C))

}
