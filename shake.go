package oak

import (
	"image"
	"math/rand"
	"time"

	"github.com/oakmound/oak/alg/floatgeom"
)

// A ScreenShaker knows how to shake a screen by a (or up to a) given magnitude.
// If Random is true, the Shaker will shake up to the (negative or positive)
// magnitude of each the X and Y axes. Otherwise, it will oscillate between
// negative magnitude and positive magnitude.
type ScreenShaker struct {
	Random    bool
	Magnitude floatgeom.Point2
}

var (
	// DefShaker is the global default shaker, used when oak.Shake is called.
	DefShaker = ScreenShaker{false, floatgeom.Point2{1.0, 1.0}}
)

// ShakeScreen will Shake using the package global DefShaker
func ShakeScreen(dur time.Duration) {
	DefShaker.Shake(dur)
}

// Shake shakes the screen based on this shaker's attributes. See ScreenShaker.
func (ss *ScreenShaker) Shake(dur time.Duration) {
	doneTime := time.Now().Add(dur)
	mag := ss.Magnitude

	setViewPos := ViewPos
	// If we end up doing this pattern more,
	// we need to replace defaultUpdateScreen
	// with a local definition of what updateScreen
	// was when this was called
	updateScreen = func(x, y int) {
		setViewPos = image.Point{x, y}
		defaultUpdateScreen(x, y)
	}
	if ss.Random {
		randOff := mag
		go func() {
			for time.Now().Before(doneTime) {
				ViewPos = setViewPos
				ViewPos.X += int(randOff.X())
				ViewPos.Y += int(randOff.Y())

				mag = mag.MulConst(-1)
				randOff = mag.MulConst(rand.Float64())
			}
			updateScreen = defaultUpdateScreen
			updateScreen(setViewPos.X, setViewPos.Y)
		}()
	} else {
		go func() {
			for time.Now().Before(doneTime) {
				ViewPos = setViewPos
				ViewPos.X += int(mag.X())
				ViewPos.Y += int(mag.Y())

				mag = mag.MulConst(-1)
			}
			updateScreen = defaultUpdateScreen
			updateScreen(setViewPos.X, setViewPos.Y)
		}()
	}
}
