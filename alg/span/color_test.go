package span

import (
	"image/color"
	"math/rand"
	"testing"
)

func TestLinearColor(t *testing.T) {
	rng := NewLinearColor(color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255})
	if rng.Poll() != (color.RGBA{255, 255, 255, 255}) {
		t.Fatal("false linear range did not return only possible value on Poll")
	}
	for i := 0; i < 100; i++ {
		if rng.Percentile(rand.Float64()) != (color.RGBA{255, 255, 255, 255}) {
			t.Fatal("false linear range did not return only possible value on Percentile")
		}
	}
	rng = NewLinearColor(color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255})
	for i := 0.0; i < 255; i++ {
		p := i / 255
		uinti := uint8(i)
		if rng.Percentile(p) != (color.RGBA{uinti, uinti, uinti, 255}) {
			t.Fatal("linear color range did not return appropriate scaled color, bottom to top")
		}
	}
	rng = NewLinearColor(color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255})
	for i := 255.0; i > 0; i-- {
		p := (255 - i) / 255
		uinti := uint8(i)
		if rng.Percentile(p) != (color.RGBA{uinti, uinti, uinti, 255}) {
			t.Fatal("linear color range did not return appropriate scaled color, top to bottom")
		}
	}
	rng = NewLinearColor(color.RGBA{125, 125, 125, 125}, color.RGBA{200, 200, 200, 200})
	if rng.Clamp(color.RGBA{100, 100, 100, 100}) != (color.RGBA{125, 125, 125, 125}) {
		t.Fatal("linear color range did not enforce minimum color")
	}
	if rng.Clamp(color.RGBA{225, 225, 225, 225}) != (color.RGBA{200, 200, 200, 200}) {
		t.Fatal("linear color range did not enforce maximum color")
	}
	if rng.Clamp(color.RGBA{175, 175, 175, 175}) != (color.RGBA{175, 175, 175, 175}) {
		t.Fatal("linear color range did not pass through value within range")
	}
}
