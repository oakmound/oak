package oak

import (
	"math/rand"
	"time"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/alg/floatgeom"
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
	// DefShaker is the global default shaker, used when oak.ShakeScreen is called.
	DefShaker = ScreenShaker{false, floatgeom.Point2{1.0, 1.0}}
)

// ShakeScreen will Shake using the package global DefShaker
func ShakeScreen(dur time.Duration) {
	DefShaker.Shake(dur)
}

// Shake shakes the screen based on this ScreenShaker's attributes.
// See DefShaker for an example shaker setup
func (ss *ScreenShaker) Shake(dur time.Duration) {
	doneTime := time.Now().Add(dur)
	mag := ss.Magnitude
	delta := intgeom.Point2{}

	if ss.Random {
		randOff := mag
		go func() {

			for time.Now().Before(doneTime) {
				xDelta := int(randOff.X())
				yDelta := int(randOff.Y())
				ShiftScreen(xDelta-delta.X(), yDelta-delta.Y())
				delta = intgeom.Point2{xDelta, yDelta}
				mag = mag.MulConst(-1)
				randOff = mag.MulConst(rand.Float64())
			}
			ShiftScreen(-delta.X(), -delta.Y())
		}()
	} else {
		go func() {

			for time.Now().Before(doneTime) {
				xDelta := int(mag.X())
				yDelta := int(mag.Y())

				ShiftScreen(xDelta, yDelta)
				delta = delta.Add(intgeom.Point2{xDelta, yDelta})
				mag = mag.MulConst(-1)
			}
			ShiftScreen(-delta.X(), -delta.Y())
		}()
	}
}
