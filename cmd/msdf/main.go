package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 256,
		Advance:    256,
	}

	C := 'R'

	tex, _ := msdf.New("/home/mo/.local/share/fonts/Hack/HackNerdFont-Regular.ttf", cfg)

	glyph := tex.Get(C)

	glyph.Save(fmt.Sprintf("%c.png", C))

}
