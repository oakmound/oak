package sequence

import "github.com/oakmound/oak/v3/audio/synth"

type WavePattern []synth.Wave

type HasWaves interface {
	GetWavePattern() []synth.Wave
	SetWavePattern([]synth.Wave)
}

func (wp *WavePattern) GetWavePattern() []synth.Wave {
	return *wp
}

func (wp *WavePattern) SetWavePattern(ws []synth.Wave) {
	*wp = ws
}

// Waves sets the generator's Wave pattern
func Waves(ws ...synth.Wave) Option {
	return func(g Generator) {
		if hw, ok := g.(HasWaves); ok {
			hw.SetWavePattern(ws)
		}
	}
}

// WaveAt sets the n'th value in the entire play sequence
// to be Wave p. This could involve duplicating a pattern
// until it is long enough to reach n. Meaningless if the
// Wave pattern has not been set yet.
func WaveAt(w synth.Wave, n int) Option {
	return func(g Generator) {
		if hw, ok := g.(HasWaves); ok {
			if hl, ok := hw.(HasLength); ok {
				if hl.GetLength() < n {
					Waves := hw.GetWavePattern()
					if len(Waves) == 0 {
						return
					}
					// If the pattern is not long enough, there are two things
					// we could do-- 1. Extend the pattern and replace the
					// individual note, or 2. Replace the note that would be
					// played at n and thus all earlier and later plays within
					// the pattern as well.
					//
					// This uses approach 1.
					for len(Waves) < n {
						Waves = append(Waves, Waves...)
					}
					Waves[n] = w
					hw.SetWavePattern(Waves)
				}
			}
		}
	}
}

// WavePatternAt sets the n'th value in the Wave pattern
// to be Wave p. Meaningless if the Wave pattern has not
// been set yet.
func WavePatternAt(w synth.Wave, n int) Option {
	return func(g Generator) {
		if hw, ok := g.(HasWaves); ok {
			Waves := hw.GetWavePattern()
			if len(Waves) < n {
				return
			}
			Waves[n] = w
			hw.SetWavePattern(Waves)
		}
	}
}
