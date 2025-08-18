package msdf

import (
	"fmt"
	"testing"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/image/math/fixed"
)

var bruteforce = &BruteForceMinDistanceFinder{
	Step: 0.01,
}

var subvision = &SubvisionMinDistanceFinder{
	FlatnessThreshold: 1.0,
	MinSegmentLength:  0.05,
	MaxSubdivisions:   1000,
	ScaleFactor:       20,
}

func TestSubvisionGetDistance(t *testing.T) {
	curve := NewCubicBezier(fixed.P(10, 10), fixed.P(10, 60), fixed.P(60, 60), fixed.P(60, 10))
	Q := Point{X: 35, Y: 20}

	T1, D1 := bruteforce.Find(curve, Q)
	T2, D2 := subvision.Find(curve, Q)

	fmt.Printf("t=%f, d=%f\n", T1, D1)
	fmt.Printf("t=%f, d=%f\n", T2, D2)

}

func TestMSDF(t *testing.T) {
	desktop, _ := homedir.Expand("~/Desktop/")
	cfg := &Config{
		Scale:          1,
		DistanceFinder: bruteforce,
		Debug:          desktop,
	}

	gen, _ := New("/home/mo/.local/share/fonts/FiraCode/FiraCodeNerdFont-Regular.ttf", cfg)
	x := gen.Get('Q')
	x.Save("/home/mo/Desktop/A.png")
	gen.Debug('9', x)

}
