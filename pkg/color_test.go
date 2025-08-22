package msdf

import "testing"

func TestSimpleEdgeColorzer(t *testing.T) {

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"smooth", func(t *testing.T) {

		}},
		{"teardrop", func(t *testing.T) {}},
		{"multicorner", func(t *testing.T) {
			t.Error("test")
		}},
	}

	for _, c := range tests {
		t.Run(c.name, c.test)
	}
}

func TestEdgeColor(t *testing.T) {

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"YELLOW", func(t *testing.T) {
			if (RED | GREEN) != YELLOW {
				t.Errorf("YELLOW is not made of  RED,GREEN")
			}
		}},
		{"MAGENTA", func(t *testing.T) {
			if (RED | BLUE) != MAGENTA {
				t.Errorf("MAGNETA is not made of  RED,BLUE")
			}
		}},
		{"CYAN", func(t *testing.T) {
			if (GREEN | BLUE) != CYAN {
				t.Errorf("CYAN is not made of  GREEN,BLUE")
			}
		}},
		{"RED_BLUE_GREEN", func(t *testing.T) {
			if RED == BLUE || BLUE == GREEN || GREEN == RED {
				t.Errorf("RED, BLUE or GREEN have the same value.")
			}
		}},
		{"RGB", func(t *testing.T) {
			c := WHITE.RGB()
			if c.R != 255 && c.B != 255 && c.G != 255 && c.A != 255 {
				t.Errorf("RGB() is not coverting correctly.")
			}
		}},
	}

	for _, c := range tests {
		t.Run(c.name, c.test)
	}
}

func TestEdgeColorPalette(t *testing.T) {

	seed := uint(0)
	palette := newEdgeColorPalette(&seed)
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"Init", func(t *testing.T) {
			c := palette.Init()
			if c != MAGENTA && c != YELLOW && c != CYAN {
				t.Error("The initial color is not MAGENTA, YELLOW or CYAN")
			}
		}},
		{"Shuffle", func(t *testing.T) {
			c := palette.Init()
			e := c
			palette.Shuffle(&c)

			if c == e {
				t.Error("Palette is not shuffling the color.")
			}
		}},
		{"ShuffleEx", func(t *testing.T) {
			c := palette.Init()
			e := c
			palette.ShuffleEx(&c, MAGENTA)

			if c == e || c == MAGENTA {
				t.Error("Palette is not shuffling the color.")
			}
		}},
	}

	for _, c := range tests {
		t.Run(c.name, c.test)
	}
}
