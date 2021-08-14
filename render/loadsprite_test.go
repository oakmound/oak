package render

import (
	"image/color"
	"testing"
)

func Test_loadSpriteNoCache_maxFileSize(t *testing.T) {
	rgba, err := loadSpriteNoCache("testdata/assets/images/16x16/jeremy.png", 1)
	if err != nil {
		t.Fatalf("failed to load jeremy: %v", err)
	}
	if rgba == nil {
		t.Fatalf("failed to load jeremy rgba")
	}
	for x := rgba.Rect.Min.X; x < rgba.Rect.Max.X; x++ {
		for y := rgba.Rect.Min.Y; y < rgba.Rect.Max.Y; y++ {
			c := rgba.RGBAAt(x, y)
			if c != (color.RGBA{0, 0, 0, 0}) {
				t.Fatal("image was not blank")
			}
		}
	}
}

func Test_loadSpriteNoCache_maxFileSize_badImage(t *testing.T) {
	_, err := loadSpriteNoCache("testdata/assets/images/16x16/bad.png", 1)
	if err == nil {
		t.Fatalf("loading bad file should have errored")
	}
}

func Test_loadSpriteNoCache_badFileExtension(t *testing.T) {
	_, err := loadSpriteNoCache("testdata/assets/images/devfile.pdn", 0)
	if err == nil {
		t.Fatalf("loading pdn file should have errored")
	}
}
