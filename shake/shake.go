// Package shake provides methods for rapidly shifting graphical components' positions
package shake

import (
	"context"
	"math/rand"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/window"
)

// A Shaker knows how to shake something by a (or up to a) given magnitude.
// If Random is true, the Shaker will shake up to the (negative or positive)
// magnitude of each the X and Y axes. Otherwise, it will oscillate between
// negative magnitude and positive magnitude.
type Shaker struct {
	Magnitude floatgeom.Point2
	Delay     time.Duration
	Random    bool
	// ResetPosition determines whether the shaken entity will be reset back to its original position
	// after shaking is complete. True by default.
	ResetPosition bool
}

var (
	// DefaultShaker is the global default shaker, used when shake.Screen or shake.Shake are called.
	DefaultShaker = &Shaker{
		Random:        false,
		Magnitude:     floatgeom.Point2{3.0, 3.0},
		Delay:         30 * time.Millisecond,
		ResetPosition: true,
	}
)

// A ShiftPoser can have its position shifted by an x,y pair
type ShiftPoser interface {
	ShiftPos(x, y float64)
}

// Shake shakes a ShiftPoser for the given duration. It uses the settings
// in DefaultShaker to determine the quality of the shake.
func Shake(sp ShiftPoser, dur time.Duration) {
	DefaultShaker.Shake(sp, dur)
}

// Shake shakes a ShiftPoser for the given duration.
func (sk *Shaker) Shake(sp ShiftPoser, dur time.Duration) {
	sk.ShakeContext(context.Background(), sp, dur)
}

// ShakeContext shakes a ShiftPoser for the given duration or until the context is done,
// whichever comes first.
func (sk *Shaker) ShakeContext(ctx context.Context, sp ShiftPoser, dur time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, dur)
	mag := sk.Magnitude
	delta := floatgeom.Point2{}

	if sk.Random {
		randOff := mag
		go func() {
			defer cancel()
			tick := time.NewTicker(sk.Delay)
			defer tick.Stop()
			for {
				select {
				case <-ctx.Done():
					if sk.ResetPosition {
						sp.ShiftPos(-delta.X(), -delta.Y())
					}
					return
				case <-tick.C:
				}
				xDelta := randOff.X() - delta.X()
				yDelta := randOff.Y() - delta.Y()
				sp.ShiftPos(xDelta, yDelta)
				delta = delta.Add(floatgeom.Point2{xDelta, yDelta})
				mag = mag.MulConst(-1)
				randOff = mag.MulConst(rand.Float64())
			}

		}()
	} else {
		go func() {
			defer cancel()
			tick := time.NewTicker(sk.Delay)
			defer tick.Stop()
			for {
				select {
				case <-ctx.Done():
					if sk.ResetPosition {
						sp.ShiftPos(-delta.X(), -delta.Y())
					}
					return
				case <-tick.C:
				}
				xDelta := mag.X()
				yDelta := mag.Y()

				sp.ShiftPos(xDelta, yDelta)
				delta = delta.Add(floatgeom.Point2{xDelta, yDelta})
				mag = mag.MulConst(-1)
			}
		}()
	}
}

type screenToPoser struct {
	window.Window
}

func (stp screenToPoser) ShiftPos(x, y float64) {
	stp.ShiftViewport(intgeom.Point2{int(x), int(y)})
}

// Screen shakes the screen that the context controls for the given duration.
// It uses the settings in DefaultShaker to determine the quality of the shake.
func Screen(ctx *scene.Context, dur time.Duration) {
	DefaultShaker.ShakeScreen(ctx, dur)
}

// ShakeScreen shakes the screen that the context controls for the given duration.
func (sk *Shaker) ShakeScreen(ctx *scene.Context, dur time.Duration) {
	poser := screenToPoser{ctx.Window}
	sk.ShakeContext(ctx, poser, dur)
}
