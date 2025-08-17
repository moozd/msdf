package msdf

import (
	"fmt"
	"math"
)

type MinDistanceFinder interface {
	GetDistance(c Curve, Q Point) (float64, float64)
}

type BruteForceMinDistanceFinder struct{}

func (s *BruteForceMinDistanceFinder) GetDistance(c Curve, Q Point) (float64, float64) {
	T := 0.0
	D := math.MaxFloat32
	for t := 0.0; t <= 1.0; t += 0.01 {
		p := c.PointAt(t)
		d := vec().fromAB(p, Q).Distance()

		if d < D {
			D = d
			T = t
		}
	}
	return T, D
}

type SubvisionMinDistanceFinder struct {
	threshold float64
}

func (s *SubvisionMinDistanceFinder) GetDistance(c Curve, Q Point) (float64, float64) {

	queue := []float64{0.0, 0.1}
	candidates := []float64{}

	for len(queue) >= 2 {
		ts := queue[0:1][0]
		te := queue[1:2][0]
		queue = queue[2:]

		cs := c.CurvatureAt(ts).asVector().Distance()
		ce := c.CurvatureAt(te).asVector().Distance()

		fmt.Printf("%f\n", math.Abs(ce-cs))

		if math.Abs(ce-cs) < 1e-8 {
			candidates = append(candidates, cs, ce)
		}
		queue = append(queue, ts, te/2.0, te/2.0, te)

	}

	fmt.Println(candidates)

	return 0.0, 0.0
}
