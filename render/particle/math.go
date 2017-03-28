package particle

import (
	"image/color"
	"math"
	"math/rand"
)

// floatFromSpread returns a random value between
// 0 and a given float64 f
func floatFromSpread(f float64) float64 {
	return (f * 2 * rand.Float64()) - f
}

// roundFloat returns a properly rounded
// integer of a given float64
func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

// randColor returns a random color from two arguments:
// a base color and a color representing the maximum
// potential offset for each of R,G,B, and A.
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

// uint16Spread returns a random uint16 between
// n-r/2 and n+r/2, not higher than 2^16-1
func uint16Spread(n, r uint32) uint16 {
	return uint16(math.Min(float64(int(n)+roundFloat(floatFromSpread(float64(r)))), 65535.0))
}

// uint16OnScale returns a uint16, progress % between n and endN.
// At 0 progress, endN will be returned. At 1 progress, n will be returned.
func uint16OnScale(n, endN uint32, progress float64) uint16 {
	return uint16((float64(n) - float64(n)*(1.0-progress) + float64(endN)*(1.0-progress)))
}
