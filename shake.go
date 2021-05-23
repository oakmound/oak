package oak

import (
	"math/rand"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
)

// TODO: Shakers don't need to be screen-dependant-- they just need something with
// a ShiftPos function.
// TODO: Shakers should accept a speed, so they aren't just moving as fast as possible

// A ScreenShaker knows how to shake a screen by a (or up to a) given magnitude.
// If Random is true, the Shaker will shake up to the (negative or positive)
// magnitude of each the X and Y axes. Otherwise, it will oscillate between
// negative magnitude and positive magnitude.
type ScreenShaker struct {
	Random    bool
	Magnitude floatgeom.Point2
}

var (
	// DefaultShaker is the global default shaker, used when ShakeScreen is called.
	DefaultShaker = &ScreenShaker{false, floatgeom.Point2{1.0, 1.0}}
)

// ShakeScreen will Shake using the package global DefaultShaker
func (c *Controller) ShakeScreen(dur time.Duration) {
	c.Shake(DefaultShaker, dur)
}

// Shake shakes the screen based on this ScreenShaker's attributes.
// See DefaultShaker for an example shaker setup
func (c *Controller) Shake(ss *ScreenShaker, dur time.Duration) {
	doneTime := time.Now().Add(dur)
	mag := ss.Magnitude
	delta := intgeom.Point2{}

	if ss.Random {
		randOff := mag
		go func() {

			for time.Now().Before(doneTime) {
				xDelta := int(randOff.X())
				yDelta := int(randOff.Y())
				c.ShiftScreen(xDelta-delta.X(), yDelta-delta.Y())
				delta = intgeom.Point2{xDelta, yDelta}
				mag = mag.MulConst(-1)
				randOff = mag.MulConst(rand.Float64())
			}
			c.ShiftScreen(-delta.X(), -delta.Y())
		}()
	} else {
		go func() {

			for time.Now().Before(doneTime) {
				xDelta := int(mag.X())
				yDelta := int(mag.Y())

				c.ShiftScreen(xDelta, yDelta)
				delta = delta.Add(intgeom.Point2{xDelta, yDelta})
				mag = mag.MulConst(-1)
			}
			c.ShiftScreen(-delta.X(), -delta.Y())
		}()
	}
}
