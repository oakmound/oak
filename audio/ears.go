//+build windows

package audio

import (
	"math"

	"bitbucket.org/StephenPatrick/go-winaudio/winaudio"
)

var (
	usingEars        bool
	earX             *float64
	earY             *float64
	earPanWidth      float64
	earSilenceRadius float64
)

// For Pan and volume calculation
func SetEars(x, y *float64, panWidth float64, silentRadius float64) {
	earX = x
	earY = y
	earPanWidth = panWidth
	earSilenceRadius = silentRadius
	usingEars = true
}

func CalculatePan(x2 float64) int32 {
	v := (x2 - *earX) * (winaudio.RIGHT_PAN / earPanWidth)
	if v < winaudio.LEFT_PAN {
		return winaudio.LEFT_PAN
	} else if v > winaudio.RIGHT_PAN {
		return winaudio.RIGHT_PAN
	}
	return int32(v)
}

func CalculateVolume(x2, y2 float64) int32 {
	// This and pan both assume a linear scale
	return int32(pointDistance(*earX, *earY, x2, y2) * (winaudio.MIN_VOLUME / earSilenceRadius))
}

func pointDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}
