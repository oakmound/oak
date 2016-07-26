package particle

import (
	"image/color"
	"math"
	"math/rand"
)

func floatFromSpread(f float64) float64 {
	return (f * 2 * rand.Float64()) - f
}

func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

func randColor(c, ra color.Color) color.Color {
	r, g, b, a := c.RGBA()
	r2, g2, b2, a2 := ra.RGBA()
	return color.RGBA64{
		uint16Spread(r, r2),
		uint16Spread(g, g2),
		uint16Spread(b, b2),
		uint16Spread(a, a2),
	}
}

func uint16Spread(n, r uint32) uint16 {
	return uint16(math.Min(float64(int(n)+roundFloat(floatFromSpread(float64(r)))), 65535.0))
}

func uint8Spread(n, r uint32) uint8 {
	n = n / 257
	r = r / 257
	return uint8(math.Min(float64(int(n)+roundFloat(floatFromSpread(float64(r)))), 255.0))
}

func unit16OnScale(n, endN uint32, progress float64) uint16 {
	return uint16((float64(n) - float64(n)*(1.0-progress) + float64(endN)*(1.0-progress)))
}

func unit8OnScale(n, endN uint32, progress float64) uint8 {
	return uint8((float64(n) - float64(n)*(1.0-progress) + float64(endN)*(1.0-progress)) / 257)
}
