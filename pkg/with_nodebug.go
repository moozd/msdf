//go:build !debug

package msdf

func (sc SimpleEdgeColorizer) verbose(c string, contour *Contour) {
}

func logf(format string, args ...interface{}) {
}
