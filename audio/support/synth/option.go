package synth

import "time"

// Option types modify waveform sources before they generate a waveform
type Option func(Source) Source

// Duration sets the duration of a generated waveform
func Duration(t time.Duration) Option {
	return func(s Source) Source {
		s.Seconds = t.Seconds()
		return s
	}
}

// Volume sets the volume of a generated waveform. It guarantees that 0 <= v <= 1
// (silent <= v <= max volume)
func Volume(v float64) Option {
	return func(s Source) Source {
		if v > 1.0 {
			v = 1.0
		} else if v < 0 {
			v = 0
		}
		s.Volume = v
		return s
	}
}

// AtPitch sets the pitch of a generated waveform.
func AtPitch(p Pitch) Option {
	return func(s Source) Source {
		s.Pitch = p
		return s
	}
}

// Mono sets the format to play mono audio.
func Mono() Option {
	return func(s Source) Source {
		s.Channels = 1
		return s
	}
}

// Stereo sets the format to play stereo audio.
func Stereo() Option {
	return func(s Source) Source {
		s.Channels = 2
		return s
	}
}
