package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	msdf "github.com/moozd/msdf/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msdf",
	Short: "Msdf texture generator",
	Long:  `This is a go implementation of msdf texture generation. check this https://github.com/Chlumsky/msdfgen `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MSDF cli use --help for more information")
	},
}

func init() {

	var glyphCmd = &cobra.Command{
		Use:   "glyph",
		Short: "Create a msdf glyph",
		Long:  "It will generate a new msdf glyph",
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := cmd.Flags().GetString("font")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			output, err := cmd.Flags().GetString("out")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			c, err := cmd.Flags().GetString("char")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			char := []rune(c)[0]
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
				Debug: debugPath,
			}
			msdfgen, _ := msdf.New(fontFile, cfg)
			s := msdfgen.Get(char)

			s.Save(filepath.Join(outDir, fmt.Sprintf("%c.png", char)))

		},
	}
	glyphCmd.Flags().StringP("font", "f", "", "Font path.")
	glyphCmd.Flags().StringP("char", "c", "", "Character.")
	glyphCmd.Flags().StringP("out", "o", ".", "Output dir path.")
	glyphCmd.Flags().BoolP("debug", "d", false, "Generate Debug output to see the edge coloring")

	rootCmd.AddCommand(glyphCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
