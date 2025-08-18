package msdf

import (
	"math"
)

type MinDistanceFinder interface {
	Find(c Curve, Q Point) (float64, float64)
}

type BruteForceMinDistanceFinder struct {
	Step float64
}

func (s *BruteForceMinDistanceFinder) Find(c Curve, Q Point) (float64, float64) {
	T := 0.0
	D := math.MaxFloat32
	step := s.Step
	if step == 0 {
		step = 0.01
	}
	for t := 0.0; t <= 1.0; t += step {
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
	MaxSubdivisions   int
	MinSegmentLength  float64
	FlatnessThreshold float64
	ScaleFactor       int
}

func (s *SubvisionMinDistanceFinder) Find(c Curve, Q Point) (float64, float64) {

	queue := []float64{0.0, 1}
	T := 0.0
	D := math.MaxFloat64

	for i := 1; len(queue) >= 2 && i < s.MaxSubdivisions; i += 1 {

		ts := queue[0]
		te := queue[1]
		queue = queue[2:]

		distance := c.PointAt(te).Sub(c.PointAt(ts)).Distance()
		flatness := c.CurvatureAt(te).Sub(c.CurvatureAt(ts)).Distance()

		if flatness < s.FlatnessThreshold || distance < s.MinSegmentLength {

			start := c.PointAt(ts)
			end := c.PointAt(te)
			aq := Q.Sub(start)
			ae := end.Sub(start)

			if aq.Distance() == 0 {
				continue
			}

			d := math.Abs(ae.Cross(aq)) / ae.Distance()

			if d < D {
				D = d
				projection := aq.Dot(ae) / ae.Dot(ae)
				T = clamp(ts+projection*(te-ts), ts, te)
			}
			continue
		}

		step := (te - ts) / float64(s.ScaleFactor)
		for j := 0; j < s.ScaleFactor; j++ {
			t_start := ts + float64(j)*step
			t_end := ts + float64(j+1)*step
			queue = append(queue, t_start, t_end)
		}

	}

	return T, D
}
