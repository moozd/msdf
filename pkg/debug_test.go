package msdf_test

import(
	"github.com/moozd/msdf/pkg"
	"testing"
)

func TestMsdf_Render(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		addr string
		cfg  *msdf.Config
		// Named input parameters for target function.
		path string
		tex  *msdf.Glyph
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := msdf.New(tt.addr, tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			m.Debug(tt.path, tt.tex)
		})
	}
}

