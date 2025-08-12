package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		Padding: 0.1,
		Debug:   true,
	}

	bank := []rune{'A', 'x', '+', '='}

	tex, _ := msdf.New("/Users/mohammad.mohammadzade/Library/Fonts/FiraCodeNerdFont-Regular.ttf", cfg)
	for _, c := range bank {
		fmt.Println(string(c))
		glyph := tex.Get(c)
		glyph.Save(fmt.Sprintf("assets/%c.png", c))
		fmt.Println()
	}

}
