package shake

import (
	"context"
	"math/rand"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/window"
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
	// DefaultShaker is the global default shaker, used when ShakeScreen is called.
	DefaultShaker = &Shaker{
		Random:        false,
		Magnitude:     floatgeom.Point2{3.0, 3.0},
		Delay:         30 * time.Millisecond,
		ResetPosition: true,
	}
)

type ShiftPoser interface {
	ShiftPos(x, y float64)
	SetPos(x, y float64)
}

func Shake(sp ShiftPoser, dur time.Duration) {
	DefaultShaker.Shake(sp, dur)
}

func (sk *Shaker) Shake(sp ShiftPoser, dur time.Duration) {
	sk.ShakeContext(context.Background(), sp, dur)
}

func (sk *Shaker) ShakeContext(ctx context.Context, sp ShiftPoser, dur time.Duration) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(dur))
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
	stp.ShiftScreen(int(x), int(y))
}

func (stp screenToPoser) SetPos(x, y float64) {
	stp.SetScreen(int(x), int(y))
}

func ShakeScreen(ctx *scene.Context, dur time.Duration) {
	DefaultShaker.ShakeScreen(ctx, dur)
}

func (sk *Shaker) ShakeScreen(ctx *scene.Context, dur time.Duration) {
	poser := screenToPoser{ctx.Window}
	sk.ShakeContext(ctx, poser, dur)
}
