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
	if v > 1.0 {
		v = 1.0
	} else if v < 0 {
		v = 0
	}
	return func(s Source) Source {
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

// Mono sets a synth source to play mono audio.
func Mono() Option {
	return func(s Source) Source {
		s.Channels = 1
		return s
	}
}

// Stereo sets a synth source to play stereo audio.
func Stereo() Option {
	return func(s Source) Source {
		s.Channels = 2
		return s
	}
}

// Detune detunes between -1.0 and 1.0, 1.0 representing a half step up.
// Q: What is detuning? A: It's taking the pitch of the audio and adjusting it less than
// a single tone up or down. If you detune too far, you've just made the next pitch,
// but if you detune a little, you get a resonant sound.
func Detune(percent float64) Option {
	return func(src Source) Source {
		curPitch := src.Pitch
		var nextPitch Pitch
		if percent > 0 {
			nextPitch = curPitch.Up(HalfStep)
		} else {
			nextPitch = curPitch.Down(HalfStep)
		}
		rawDelta := float64(int16(curPitch) - int16(nextPitch))
		delta := rawDelta * percent
		// TODO: does pitch need to be a float?
		src.Pitch = Pitch(float64(curPitch) + delta)
		return src
	}
}
