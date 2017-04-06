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
	fps := 1 / now.Sub(lastTime).Seconds()
	// This indicates that time.Now recorded two times within
	// the innacuracy of the OS's system clock, so the values
	// were the same.
	if int(fps) < 0 {
		return 1200
	}
	return fps
}

// FPSToNano converts a framesPerSecond value to the number of
// nanoseconds that should take place for each frame.
func FPSToNano(fps float64) int64 {
	return int64(nanoPerSecond / fps)
}
