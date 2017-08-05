package render

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugIdentity(t *testing.T) {
	r := NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	UpdateDebugMap("r", r)
	r2, ok := GetDebugRenderable("r")
	assert.True(t, ok)
	assert.Equal(t, r, r2)
}
