//go:build debug

package msdf

import "log"

func (sc SimpleEdgeColorizer) verbose(c string, contour *Contour) {
	log.Printf("Colorizer: case=%s char=%c\n", c, contour.Symbol)
	for _, e := range contour.Edges {
		log.Printf("Colorizer: edge=%v %v\n", e, e.Curve)
	}
}

func logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
