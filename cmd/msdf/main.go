package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		LineHeight: 250,
		Advance:    250,
		Padding:    0.5,
		Debug:      true,
	}

	bank := []rune{'A'}

	tex, _ := msdf.New("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf", cfg)
	for _, c := range bank {
		fmt.Println(string(c))
		glyph := tex.Get(c)
		glyph.Save(fmt.Sprintf("assets/%c.png", c))
		fmt.Println()
	}

}
