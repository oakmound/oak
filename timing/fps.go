package timing

import (
	"math"
	"time"
)

var (
	nanoPerSecond = math.Pow(10, 9)
)

// FPS returns the number of frames being processed per second,
// supposing a time interval from lastTime to now.
func FPS(lastTime, now time.Time) float64 {
	return 1 / float64(now.Sub(lastTime).Seconds())
}

// FPStoNano converts a framesPerSecond value to the number of
// nanoseconds that should take place for each frame.
func FPSToNano(fps float64) int64 {
	return int64(nanoPerSecond / fps)
}
