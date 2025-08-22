package msdf

import (
	"testing"

	"golang.org/x/image/math/fixed"
)

func TestSign(t *testing.T) {
	z := sign(0.0)
	n := sign(-100.0)
	p := sign(+100.0)

	if z != 0 {
		t.Errorf("Expected 0 but got %f", z)
	}

	if n >= 0.0 {
		t.Errorf("Expected -1 but got %f", n)
	}

	if p <= 0.0 {
		t.Errorf("Expected +1 but got %f", p)
	}
}

func TestLerp(t *testing.T) {

	p0 := fixed.P(0, 0)
	p1 := fixed.P(10, 10)
	p := lerp(p0, p1, 0.5)

	if p.X != fixed.I(5) {
		t.Errorf("Expected x=%v but got x=%v", fixed.I(5), p.X)
	}

	if p.Y != fixed.I(5) {
		t.Errorf("Expected y=%v but got y=%v", fixed.I(5), p.Y)
	}
}

func TestClamp(t *testing.T) {
	a := clamp(3.0, 0, 1)
	b := clamp(0.5, 0, 1)
	c := clamp(-0.1, 0, 1)

	if a != 1.0 {
		t.Errorf("Expected %f but got %f", 1.0, a)
	}

	if b != 0.5 {
		t.Errorf("Expected %f but got %f", 0.5, b)
	}

	if c != 0.0 {
		t.Errorf("Expected %f but got %f", 0.0, c)
	}
}
func TestPackI26_6(t *testing.T) {
	a := packI26_6(5.0 + (1 / 64.0))   // 26_6 -> 6 bits -> 2^6=64
	e := fixed.I(5) + fixed.Int26_6(1) // Int26_6(1) == 1/64
	if a != e {
		t.Errorf("Expected %v but got %v", e, a)
	}

}

func TestUnpackI26_6(t *testing.T) {

	a := fixed.I(5)       //5
	b := fixed.Int26_6(1) // 1/64

	e1 := 5.0
	e2 := 1 / 64.0
	t1 := unpackI26_6(a)
	t2 := unpackI26_6(b)

	if t1 != e1 {
		t.Errorf("Expected %f but got %f", e1, t1)
	}

	if t2 != e2 {
		t.Errorf("Expected %f but got %f", e2, t2)
	}
}

func TestPackP26_6(t *testing.T) {
	p := packP26_6(3.0, 3.0)
	e := fixed.I(3)

	if p.X != e {
		t.Errorf("Expected %v but got %v", e, p.X)
	}

	if p.Y != e {
		t.Errorf("Expected %v but got %v", e, p.Y)
	}
}

func TestUnpackP26_6(t *testing.T) {
	x, y := unpackP26_6(fixed.P(3, 3))
	e := 3.0

	if x != e {
		t.Errorf("Expected %f but got %f", e, x)
	}

	if y != e {
		t.Errorf("Expected %f but got %f", e, y)
	}

}
