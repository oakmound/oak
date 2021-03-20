package render

import (
	"image"
	"image/color"
	"math"
	"reflect"
	"testing"
)

func TestLine(t *testing.T) {
	l := NewLine(0, 0, 10, 10, color.RGBA{255, 255, 255, 255})
	rgba := l.GetRGBA()
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if x == y {
				if rgba.At(x, y) != (color.RGBA{255, 255, 255, 255}) {
					t.Fatalf("rgba pixel mismatch")
				}
			} else {
				if rgba.At(x, y) != (color.RGBA{0, 0, 0, 0}) {
					t.Fatalf("rgba pixel mismatch")
				}
			}
		}
	}
	l = NewLine(0, 0, 0, 0, color.RGBA{255, 255, 255, 255})
	rgba = l.GetRGBA()
	rgba2 := image.NewRGBA(image.Rect(0, 0, 1, 1))
	rgba2.Set(0, 0, color.RGBA{255, 255, 255, 255})
	if !reflect.DeepEqual(rgba, rgba2) {
		t.Fatalf("manually drawn dot/line did not match new dot/line")
	}

	l = NewLine(0, 0, 0, 5, color.RGBA{255, 255, 255, 255})
	rgba = l.GetRGBA()
	rgba2 = image.NewRGBA(image.Rect(0, 0, 1, 5))
	for y := 0; y < 5; y++ {
		rgba2.Set(0, y, color.RGBA{255, 255, 255, 255})
	}
	if !reflect.DeepEqual(rgba, rgba2) {
		t.Fatalf("manually drawn line did not match new line")
	}
}

func TestThickLine(t *testing.T) {
	l := NewThickLine(0, 0, 10, 10, color.RGBA{255, 255, 255, 255}, 1)
	rgba := l.GetRGBA()
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if math.Abs(float64(x)-float64(y)) <= 2 {
				if rgba.At(x, y) != (color.RGBA{255, 255, 255, 255}) {
					t.Fatalf("rgba pixel mismatch")
				}
			} else {
				if rgba.At(x, y) != (color.RGBA{0, 0, 0, 0}) {
					t.Fatalf("rgba pixel mismatch")
				}
			}
		}
	}
}

//TODO: Update to use progress function to test coloring
func TestGradientLine(t *testing.T) {
	l := NewGradientLine(0, 0, 10, 10, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, 1)
	rgba := l.GetRGBA()
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if math.Abs(float64(x)-float64(y)) <= 2 {
				if rgba.At(x, y) != (color.RGBA{255, 255, 255, 255}) {
					t.Fatalf("rgba pixel mismatch")
				}
			} else {
				if rgba.At(x, y) != (color.RGBA{0, 0, 0, 0}) {
					t.Fatalf("rgba pixel mismatch")
				}
			}
		}
	}
}

func TestDrawLine(t *testing.T) {
	l := NewLine(0, 0, 10, 10, color.RGBA{255, 255, 255, 255})
	rgba := l.GetRGBA()
	// See height addition in line
	rgba2 := image.NewRGBA(image.Rect(0, 0, 10, 11))
	DrawLine(rgba2, 0, 0, 10, 10, color.RGBA{255, 255, 255, 255})
	if !reflect.DeepEqual(rgba, rgba2) {
		t.Fatalf("draw line did not match new line")
	}
	rgba3 := image.NewRGBA(image.Rect(0, 0, 10, 11))
	DrawGradientLine(rgba3, 10, 10, 0, 0, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, 0)
	if !reflect.DeepEqual(rgba, rgba3) {
		t.Fatalf("gradient line from black to black did not match solid line")
	}
}

func TestThickLinePoint(t *testing.T) {
	// p1 = p2
	l := NewThickLine(0, 0, 0, 0, color.RGBA{255, 0, 0, 255}, 4)
	rgba := l.GetRGBA()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if rgba.At(i, j) != (color.RGBA{255, 0, 0, 255}) {
				t.Fatalf("rgba pixel mismatch")
			}
		}
	}
}
func TestThickLineVert(t *testing.T) {
	// Vertical
	l := NewThickLine(0, 0, 0, 10, color.RGBA{255, 0, 0, 255}, 4)
	rgba := l.GetRGBA()
	for i := 0; i < 5; i++ {
		for j := 0; j < 18; j++ {
			if rgba.At(i, j) != (color.RGBA{255, 0, 0, 255}) {
				t.Fatalf("rgba pixel mismatch")
			}
		}
	}
}
