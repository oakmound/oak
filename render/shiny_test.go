package render

import (
	"image"
	"image/color"
	"testing"
)

func TestShinyDrawFns(t *testing.T) {
	world := image.NewRGBA(image.Rect(0, 0, 20, 20))
	rgba := image.NewRGBA(image.Rect(0, 0, 20, 20))

	world.SetRGBA(10, 10, color.RGBA{255, 0, 0, 255})

	DrawImage(world, rgba, 0, 0)
	if world.At(10, 10) != (color.RGBA{255, 0, 0, 255}) {
		t.Fatalf("draw image overwrote rgba")
	}
	OverwriteImage(world, rgba, 0, 0)
	if world.At(10, 10) != (color.RGBA{0, 0, 0, 0}) {
		t.Fatalf("overwrite image did not overwrite rgba")
	}
}
