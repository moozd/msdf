package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	msdf "github.com/pierrec/msdf/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msdf",
	Short: "Msdf texture generator",
	Long:  `This is a go implementation of msdf texture generation. check this https://github.com/Chlumsky/msdfgen `,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	addr, err := cmd.Flags().GetString("font")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	charset, err := cmd.Flags().GetString("charset")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	seed, err := cmd.Flags().GetUint("seed")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scale, err := cmd.Flags().GetFloat64("scale")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	distanceField, err := cmd.Flags().GetFloat64("distance-field")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	size, err := cmd.Flags().GetFloat64("size")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fontFile, err := homedir.Expand(addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outDir, err := homedir.Expand(output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	debugPath := ""
	if debug {
		debugPath = outDir
	}
	cfg := &msdf.Config{
		Seed:             seed,
		Scale:            scale,
		Size:             size,
		DebugArtifactDir: debugPath,
		DistanceField:    distanceField,
		EdgeColorizer:    &msdf.SimpleEdgeColorizer{},
		DistanceFinder:   &msdf.BruteForceMinDistanceFinder{},
	}
	msdfgen, err := msdf.New(fontFile, cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c, m, err := msdfgen.CreateAtlas(512, charset)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.Save(path.Join(outDir, "atlas.png"))
	m.Save(path.Join(outDir, "atlas.json"))
}

func init() {

	rootCmd.Flags().BoolP("debug", "d", false, "Generate Debug output to see the edge coloring")
	rootCmd.Flags().StringP("font", "f", "", "Font path.")
	rootCmd.Flags().StringP("charset", "c", "", "Character set.")
	rootCmd.Flags().StringP("output", "o", ".", "Output dir path.")
	rootCmd.Flags().Uint("seed", 0, "coloring seed")
	rootCmd.Flags().Float64("scale", 1.0, "texture scale")
	rootCmd.Flags().Float64P("size", "s", 10.0, "font size")
	rootCmd.Flags().Float64("distance-field", 4.0, "Distance field, default is 4.0")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
