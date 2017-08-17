package render

import "math"

// Progress functions return some float64 between 0 and 1 for how far along
// some gradient the position (x,y) is on a rectangle of (w,h) dimensions
type Progress func(x, y, w, h int) float64

// ProgressAnd will combine multiple progress functions into one, such that
// each progress function contributes an equal amount of progress to the total.
// Example:
// ProgressAnd(HorizontalProgress,VerticalProgress) will return .5 for y = 0,
// x = w, .5 for y = h, x = 0, and 1.0 for y = h, x = w.
func ProgressAnd(pfns ...Progress) Progress {
	return func(x, y, w, h int) float64 {
		p := 0.0
		for _, pf := range pfns {
			p += pf(x, y, w, h)
		}
		return p / float64(len(pfns))
	}
}

// Progress functions
var (
	// HorizontalProgress measures progress as x / w
	HorizontalProgress = func(x, y, w, h int) float64 {
		return float64(x) / float64(w)
	}
	// VerticalProgress measures progress as y / h
	VerticalProgress = func(x, y, w, h int) float64 {
		return float64(y) / float64(h)
	}
	// DiagonalProgress measures progress along the x,y diagonal from upleft
	// to bottom right.
	DiagonalProgress = ProgressAnd(HorizontalProgress, VerticalProgress)
	// CircularProgress measures progress as distance from the center of a circle.
	CircularProgress = func(x, y, w, h int) float64 {
		xRadius := float64(w) / 2
		yRadius := float64(h) / 2
		dX := math.Abs(float64(x) - xRadius)
		dY := math.Abs(float64(y) - yRadius)
		progress := math.Pow(dX/xRadius, 2) + math.Pow(dY/yRadius, 2)
		if progress > 1 {
			progress = 1
		}
		return progress
	}
)
