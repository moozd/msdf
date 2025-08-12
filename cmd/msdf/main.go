package main

import (
	"fmt"

	msdf "github.com/moozd/msdf/pkg"
)

func main() {

	cfg := &msdf.Config{
		Scale: 1,
		Debug: true,
	}

	bank := []rune{'X'}

	msdfgen, _ := msdf.New("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf", cfg)
	for _, c := range bank {
		fmt.Println(string(c))
		o := msdfgen.Get(c)
		o.Save(fmt.Sprintf("assets/%c.png", c))
		msdfgen.Debug(fmt.Sprintf("assets/%c_render.png", c), o)

		fmt.Println()
	}

}
