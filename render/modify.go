package render

import (
	"image"
	"image/color"
	"math"
)

type Modifiable interface {
	FlipX()
	FlipY()
	ApplyColor(c color.Color)
	ApplyMask(img image.RGBA)
	Rotate(degrees int)
	Scale(xRatio float64, yRatio float64)
}

func FlipX(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(w-x, y))
		}
	}
	return newRgba
}

func FlipY(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, h-y))
		}
	}
	return newRgba
}

func ApplyColor(rgba *image.RGBA, c color.Color) *image.RGBA {
	r1, g1, b1, _ := c.RGBA()
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()
			tmp := color.RGBA{uint8(r1 * r2 / 255),
				uint8(g1 * g2 / 255),
				uint8(b1 * b2 / 255),
				uint8(a2)}
			newRgba.Set(x, y, tmp)
		}
	}
	return newRgba
}

func ApplyMask(rgba *image.RGBA, img image.RGBA) *image.RGBA {
	// Instead of static color it just two buffers melding
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, y))
		}
	}
	bounds = img.Bounds()
	w = bounds.Max.X
	h = bounds.Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()
			tmp := color.RGBA{uint8(r1 * r2 / 255),
				uint8(g1 * g2 / 255),
				uint8(b1 * b2 / 255),
				uint8(a2)}
			newRgba.Set(x, y, tmp)
		}
	}
	return newRgba
}

func Rotate(rgba *image.RGBA, degrees int) *image.RGBA {
	//otates clockwise by the given degrees
	// Will shear any pixels that land outside the given buffer
	angle := float64(degrees) / 180 * math.Pi
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	centerX := float64(w / 2)
	centerY := float64(h / 2)
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)

	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xf := float64(x)
			yf := float64(y)
			newX := round(cosAngle*(xf-centerX) - sinAngle*(yf-centerY) + centerX)
			newY := round(sinAngle*(xf-centerX) + cosAngle*(yf-centerY) + centerY)
			newRgba.Set(newX, newY, rgba.At(x, y))
		}
	}
	return newRgba
}

func Scale(rgba *image.RGBA, xRatio float64, yRatio float64) *image.RGBA {
	bounds := rgba.Bounds()
	w := int(math.Floor(float64(bounds.Max.X) * xRatio))
	h := int(math.Floor(float64(bounds.Max.Y) * yRatio))
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(int(math.Floor(float64(x)/xRatio)), int(math.Floor(float64(y)/yRatio))))
		}
	}
	return newRgba
}

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}
