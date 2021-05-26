package shake

import (
	"math/rand"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
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
		Magnitude:     floatgeom.Point2{1.0, 1.0},
		Delay:         30 * time.Millisecond,
		ResetPosition: true,
	}
)

type ShiftPoser interface {
	ShiftPos(x, y float64)
	SetPos(x, y float64)
	GetPos() (x, y float64)
}

func Shake(sp ShiftPoser, dur time.Duration) {
	DefaultShaker.Shake(sp, dur)
}

func (sk *Shaker) Shake(sp ShiftPoser, dur time.Duration) {
	doneTime := time.Now().Add(dur)
	mag := sk.Magnitude
	delta := floatgeom.Point2{}

	origX, origY := sp.GetPos()

	if sk.Random {
		randOff := mag
		go func() {
			tick := time.NewTicker(sk.Delay)
			defer tick.Stop()
			for {
				<-tick.C
				if time.Now().After(doneTime) {
					break
				}
				xDelta := randOff.X()
				yDelta := randOff.Y()
				sp.ShiftPos(xDelta-delta.X(), yDelta-delta.Y())
				delta = floatgeom.Point2{xDelta, yDelta}
				mag = mag.MulConst(-1)
				randOff = mag.MulConst(rand.Float64())
			}
			if sk.ResetPosition {
				sp.SetPos(origX, origY)
			}
		}()
	} else {
		go func() {
			tick := time.NewTicker(sk.Delay)
			defer tick.Stop()
			for {
				<-tick.C
				if time.Now().After(doneTime) {
					break
				}
				xDelta := mag.X()
				yDelta := mag.Y()

				sp.ShiftPos(xDelta, yDelta)
				delta = delta.Add(floatgeom.Point2{xDelta, yDelta})
				mag = mag.MulConst(-1)
			}
			if sk.ResetPosition {
				sp.SetPos(origX, origY)
			}
		}()
	}
}

type ShiftScreener interface {
	ShiftScreen(x, y int)
	SetScreen(x, y int)
	Viewport() intgeom.Point2
}

type screenToPoser struct {
	ShiftScreener
}

func (stp screenToPoser) ShiftPos(x, y float64) {
	stp.ShiftScreen(int(x), int(y))
}

func (stp screenToPoser) SetPos(x, y float64) {
	stp.SetScreen(int(x), int(y))
}

func (stp screenToPoser) GetPos() (x, y float64) {
	vp := stp.Viewport()
	return float64(vp.X()), float64(vp.Y())
}

func ShakeScreen(ss ShiftScreener, dur time.Duration) {
	DefaultShaker.ShakeScreen(ss, dur)
}

func (sk *Shaker) ShakeScreen(ss ShiftScreener, dur time.Duration) {
	poser := screenToPoser{ss}
	sk.Shake(poser, dur)
}
