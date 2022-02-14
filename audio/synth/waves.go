// Package synth provides functions and types to support waveform synthesis.
package synth

import (
	"math"

	audio "github.com/oakmound/oak/v3/audio/klang"
	"github.com/oakmound/oak/v3/audio/pcm"
)

// Wave functions take a set of options and return an audio
type Wave func(opts ...Option) (audio.Audio, error)

// Sourced from https://en.wikibooks.org/wiki/Sound_Synthesis_Theory/Oscillators_and_Wavetables
func phase(freq Pitch, i int, sampleRate uint32) float64 {
	return float64(freq) * (float64(i) / float64(sampleRate)) * 2 * math.Pi
}

func bytesFromInts(is []int16, channels int) []byte {
	wave := make([]byte, len(is)*channels*2)
	for i := 0; i < len(wave); i += channels * 2 {
		wave[i] = byte(is[i/4] % 256)
		wave[i+1] = byte(is[i/4] >> 8)
		// duplicate the contents across all channels
		for c := 1; c < channels; c++ {
			wave[i+(2*c)] = wave[i]
			wave[i+(2*c)+1] = wave[i+1]
		}
	}
	wave = append(wave, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	return wave
}

// Sin produces a Sin wave
//         __
//       --  --
//      /      \
//--__--        --__--
func (s Source) Sin(opts ...Option) (audio.Audio, error) {
	s = s.Update(opts...)
	var b []byte
	switch s.Bits {
	case 16:
		s.Volume *= 65535 / 2
		wave := make([]int16, int(s.Seconds*float64(s.SampleRate)))
		for i := 0; i < len(wave); i++ {
			wave[i] = int16(s.Volume * math.Sin(s.Phase(i)))
		}
		b = bytesFromInts(wave, int(s.Channels))
	}
	return s.Wave(b)
}

// SinPCM acts like Sin, but returns a PCM type instead of a klang type.
func (s Source) SinPCM(opts ...Option) (pcm.Reader, error) {
	s.Volume *= 65535 / 2
	return &WaveReader{
		Source: s.Update(opts...),
		waveFunc: func(s Source, idx int) int16 {
			return int16(s.Volume * math.Sin(s.Phase(idx)))
		},
	}, nil
}

// Pulse acts like Square when given a pulse of 2, when given any lesser
// pulse the time up and down will change so that 1/pulse time the wave will
// be up.
//
//     __    __
//     ||    ||
// ____||____||____
func (s Source) Pulse(pulse float64) Wave {
	pulseSwitch := 1 - 2/pulse
	return func(opts ...Option) (audio.Audio, error) {
		s = s.Update(opts...)

		var b []byte
		switch s.Bits {
		case 16:
			wave := make([]int16, int(s.Seconds*float64(s.SampleRate)))
			for i := range wave {
				// alternatively phase % 2pi
				if math.Sin(s.Phase(i)) > pulseSwitch {
					wave[i] = int16(s.Volume)
				} else {
					wave[i] = int16(-s.Volume)
				}
			}
			b = bytesFromInts(wave, int(s.Channels))
		}
		return s.Wave(b)
	}
}

type WaveReader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int16
}

func (pr *WaveReader) ReadPCM(b []byte) (n int, err error) {
	switch pr.Bits {
	case 16:
		bytesPerI16 := int(pr.Channels) * 2
		for i := 0; i+bytesPerI16 <= len(b); i += bytesPerI16 {
			i16 := pr.waveFunc(pr.Source, pr.lastIndex)
			pr.lastIndex++
			for c := 0; c < int(pr.Channels); c++ {
				b[i+(2*c)] = byte(i16)
				b[i+(2*c)+1] = byte(i16 >> 8)
			}
			n += bytesPerI16
		}
	}
	return
}

func (pr *WaveReader) PCMFormat() pcm.Format {
	return pcm.Format{
		SampleRate: pr.SampleRate,
		Channels:   pr.Channels,
		Bits:       pr.Bits,
	}
}

// PulsePCM acts like Pulse, but returns a PCM type instead of a klang type.
func (s Source) PulsePCM(pulse float64) func(opts ...Option) (pcm.Reader, error) {
	s.Volume *= 65535 / 2
	pulseSwitch := 1 - 2/pulse
	return func(opts ...Option) (pcm.Reader, error) {
		return &WaveReader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int16 {
				if math.Sin(s.Phase(idx)) > pulseSwitch {
					return int16(s.Volume)
				} else {
					return int16(-s.Volume)
				}
			},
		}, nil
	}
}

// Square produces a Square wave
//
//       _________
//       |       |
// ______|       |________
func (s Source) Square(opts ...Option) (audio.Audio, error) {
	return s.Pulse(2)(opts...)
}

// Saw produces a saw wave
//
//   ^   ^   ^
//  / | / | /
// /  |/  |/
func (s Source) Saw(opts ...Option) (audio.Audio, error) {
	s = s.Update(opts...)

	var b []byte
	switch s.Bits {
	case 16:
		s.Volume *= 65535 / 2
		wave := make([]int16, int(s.Seconds*float64(s.SampleRate)))
		for i := range wave {
			wave[i] = int16(s.Volume - (s.Volume / math.Pi * math.Mod(s.Phase(i), 2*math.Pi)))
		}
		b = bytesFromInts(wave, int(s.Channels))
	}
	return s.Wave(b)
}

// SawPCM acts like Saw, but returns a PCM type instead of a klang type.
func (s Source) SawPCM(opts ...Option) (pcm.Reader, error) {
	s.Volume *= 65535 / 2
	return &WaveReader{
		Source: s.Update(opts...),
		waveFunc: func(s Source, idx int) int16 {
			return int16(s.Volume - (s.Volume / math.Pi * math.Mod(s.Phase(idx), 2*math.Pi)))
		},
	}, nil
}

// Triangle produces a Triangle wave
//
//   ^   ^
//  / \ / \
// v   v   v
func (s Source) Triangle(opts ...Option) (audio.Audio, error) {
	s = s.Update(opts...)

	var b []byte
	switch s.Bits {
	case 16:
		s.Volume *= 65535 / 2
		wave := make([]int16, int(s.Seconds*float64(s.SampleRate)))
		for i := range wave {
			p := math.Mod(s.Phase(i), 2*math.Pi)
			m := int16(p * (2 * s.Volume / math.Pi))
			if math.Sin(p) > 0 {
				wave[i] = int16(-s.Volume) + m
			} else {
				wave[i] = 3*int16(s.Volume) - m
			}
		}
		b = bytesFromInts(wave, int(s.Channels))
	}
	return s.Wave(b)
}

// TrianglePCM acts like Triangle, but returns a PCM type instead of a klang type.
func (s Source) TrianglePCM(opts ...Option) (pcm.Reader, error) {
	s.Volume *= 65535 / 2
	return &WaveReader{
		Source: s.Update(opts...),
		waveFunc: func(s Source, idx int) int16 {
			p := math.Mod(s.Phase(idx), 2*math.Pi)
			m := int16(p * (2 * s.Volume / math.Pi))
			if math.Sin(p) > 0 {
				return int16(-s.Volume) + m
			}
			return 3*int16(s.Volume) - m
		},
	}, nil
}

// Could have pulse triangle
