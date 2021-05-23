package render

import (
	"image/color"
	"math/rand"
	"testing"
	"time"
)

func TestIdentityColorer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	c := IdentityColorer(color.RGBA{255, 100, 100, 255})
	if c(rand.Float64()) != (color.RGBA{255, 100, 100, 255}) {
		t.Fatalf("identity colorer did not return set color")
	}
}
