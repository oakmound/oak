package synth

import (
	"time"

	"github.com/oakmound/oak/v3/audio/internal/audio"
)

// A Source stores necessary information for generating audio and waveform data
type Source struct {
	audio.Format
	Pitch   Pitch
	Volume  float64
	Seconds float64
}

// PlayLength returns the time it will take before audio generated from this
// source will stop.
func (s Source) PlayLength() time.Duration {
	return time.Duration(s.Seconds) * 1000 * time.Millisecond
}

// Phase is shorthand for phase(s.Pitch, i, s.SampleRate).
// Some sources might have custom phase functions in the future, however.
func (s Source) Phase(i int) float64 {
	return phase(s.Pitch, i, s.SampleRate)
}

// Update is shorthand for applying a set of options to a source
func (s Source) Update(opts ...Option) Source {
	for _, opt := range opts {
		s = opt(s)
	}
	return s
}

var (
	// Int16 is a default source for building 16-bit audio
	Int16 = Source{
		Format: audio.Format{
			SampleRate: 44100,
			Channels:   2,
			Bits:       16,
		},
		Pitch:   A4,
		Volume:  .25,
		Seconds: 1,
	}
)
