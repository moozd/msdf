package msdf

import (
	"fmt"
	"testing"

	"golang.org/x/image/math/fixed"
)

func TestSubvisionGetDistance(t *testing.T) {
	curve := NewCubicBezier(fixed.P(10, 10), fixed.P(10, 60), fixed.P(60, 60), fixed.P(60, 10))
	Q := Point{X: 35, Y: 20}

	bruteforce := &BruteForceMinDistanceFinder{}

	T1, D1 := bruteforce.GetDistance(curve, Q)
	fmt.Printf("t=%f, d=%f\n", T1, D1)

	subvision := &SubvisionMinDistanceFinder{
		threshold: 1e-8,
	}

	T2, D2 := subvision.GetDistance(curve, Q)
	fmt.Printf("t=%f, d=%f\n", T2, D2)

}
