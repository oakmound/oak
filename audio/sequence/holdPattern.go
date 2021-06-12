package sequence

import "time"

// A HoldPattern is a pattern that might loop on itself for how long notes
// should be held
type HoldPattern []time.Duration

// HasHolds enables generators to be built from HoldPattern and use the
// related option functions
type HasHolds interface {
	GetHoldPattern() *[]time.Duration
}

// GetHoldPattern lets composing HoldPattern satisfy HasHolds
func (hp *HoldPattern) GetHoldPattern() *HoldPattern {
	return hp
}

// Holds sets the generator's Hold pattern
func Holds(vs ...time.Duration) Option {
	return func(g Generator) {
		if hhs, ok := g.(HasHolds); ok {
			*hhs.GetHoldPattern() = vs
		}
	}
}

// HoldAt sets the n'th value in the entire play sequence
// to be Hold p. This could involve duplicating a pattern
// until it is long enough to reach n. Meaningless if the
// Hold pattern has not been set yet.
func HoldAt(t time.Duration, n int) Option {
	return func(g Generator) {
		if hhs, ok := g.(HasHolds); ok {
			if hl, ok := hhs.(HasLength); ok {
				if hl.GetLength() < n {
					hp := hhs.GetHoldPattern()
					Holds := *hp
					if len(Holds) == 0 {
						return
					}
					// If the pattern is not long enough, there are two things
					// we could do-- 1. Extend the pattern and replace the
					// individual note, or 2. Replace the note that would be
					// played at n and thus all earlier and later plays within
					// the pattern as well.
					//
					// This uses approach 1.
					for len(Holds) <= n {
						Holds = append(Holds, Holds...)
					}
					Holds[n] = t
					*hp = Holds
				}
			}
		}
	}
}

// HoldPatternAt sets the n'th value in the Hold pattern
// to be Hold p. Meaningless if the Hold pattern has not
// been set yet.
func HoldPatternAt(t time.Duration, n int) Option {
	return func(g Generator) {
		if hhs, ok := g.(HasHolds); ok {
			hp := hhs.GetHoldPattern()
			Holds := *hp
			if len(Holds) <= n {
				return
			}
			Holds[n] = t
			*hp = Holds
		}
	}
}
