package sequence

import "github.com/oakmound/oak/v3/audio/support/synth"

type PitchPattern []synth.Pitch

type HasPitches interface {
	GetPitchPattern() []synth.Pitch
	SetPitchPattern([]synth.Pitch)
}

func (pp *PitchPattern) GetPitchPattern() []synth.Pitch {
	return *pp
}

func (pp *PitchPattern) SetPitchPattern(ps []synth.Pitch) {
	*pp = ps
}

// Pitches sets the generator's pitch pattern
func Pitches(ps ...synth.Pitch) Option {
	return func(g Generator) {
		if hpp, ok := g.(HasPitches); ok {
			hpp.SetPitchPattern(ps)
		}
	}
}

// PitchAt sets the n'th value in the entire play sequence
// to be pitch p. This could involve duplicating a pattern
// until it is long enough to reach n. Meaningless if the
// pitch pattern has not been set yet.
func PitchAt(p synth.Pitch, n int) Option {
	return func(g Generator) {
		if hpp, ok := g.(HasPitches); ok {
			if hl, ok := hpp.(HasLength); ok {
				if hl.GetLength() < n {
					pitches := hpp.GetPitchPattern()
					if len(pitches) == 0 {
						return
					}
					// If the pattern is not long enough, there are two things
					// we could do-- 1. Extend the pattern and replace the
					// individual note, or 2. Replace the note that would be
					// played at n and thus all earlier and later plays within
					// the pattern as well.
					//
					// This uses approach 1.
					for len(pitches) < n {
						pitches = append(pitches, pitches...)
					}
					pitches[n] = p
					hpp.SetPitchPattern(pitches)
				}
			}
		}
	}
}

// PitchPatternAt sets the n'th value in the pitch pattern
// to be pitch p. Meaningless if the pitch pattern has not
// been set yet.
func PitchPatternAt(p synth.Pitch, n int) Option {
	return func(g Generator) {
		if hpp, ok := g.(HasPitches); ok {
			pitches := hpp.GetPitchPattern()
			if len(pitches) < n {
				return
			}
			pitches[n] = p
			hpp.SetPitchPattern(pitches)
		}
	}
}
