package render

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGradient(t *testing.T) {
	a := color.RGBA{0, 0, 0, 0}
	b := color.RGBA{255, 255, 255, 255}
	for i := uint16(0); i < 255; i++ {
		progress := float64(i) / 255.0
		gc := GradientColorAt(a, b, progress)
		v := (i * 257)
		diff := v - gc.R
		assert.True(t, diff < 2)
	}
}
