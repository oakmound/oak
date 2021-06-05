package filter

import (
	"github.com/oakmound/oak/v3/audio/internal/audio"
	"github.com/oakmound/oak/v3/audio/internal/audio/filter/supports"
)

// A SampleRate is a function that takes in uint32 SampleRates
type SampleRate func(*uint32)

// Apply checks that the given audio supports SampleRate, filters if it
// can, then returns
func (srf SampleRate) Apply(a audio.Audio) (audio.Audio, error) {
	if ssr, ok := a.(supports.SampleRate); ok {
		srf(ssr.GetSampleRate())
		return a, nil
	}
	return a, supports.NewUnsupported([]string{"SampleRate"})
}

// ModSampleRate might slow down or speed up a sample, but this will
// effect the perceived pitch of the sample. See Speed.
func ModSampleRate(mult float64) SampleRate {
	return func(sr *uint32) {
		*sr = uint32(float64(*sr) * mult)
	}
}
