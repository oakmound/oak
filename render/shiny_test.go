package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShinyDrawFns(t *testing.T) {
	world := image.NewRGBA(image.Rect(0, 0, 20, 20))
	ShinySet(world, color.RGBA{255, 0, 0, 255}, -10, -10)
	assert.Equal(t, color.RGBA{255, 0, 0, 255}, world.At(10, 10))

	rgba := image.NewRGBA(image.Rect(0, 0, 20, 20))

	ShinyDraw(world, rgba, 0, 0)
	assert.Equal(t, color.RGBA{255, 0, 0, 255}, world.At(10, 10))
	ShinyOverwrite(world, rgba, 0, 0)
	assert.Equal(t, color.RGBA{0, 0, 0, 0}, world.At(10, 10))
}
