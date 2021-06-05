package sequence

import "time"

type Tick time.Duration

type HasTicks interface {
	GetTick() time.Duration
	SetTick(time.Duration)
}

func (vp *Tick) GetTick() time.Duration {
	return time.Duration(*vp)
}

func (vp *Tick) SetTick(vs time.Duration) {
	*vp = Tick(vs)
}

// Ticks sets the generator's Tick
func Ticks(t time.Duration) Option {
	return func(g Generator) {
		if ht, ok := g.(HasTicks); ok {
			ht.SetTick(t)
		}
	}
}
