// Package synth provides functions and types to support waveform synthesis.
package synth

import (
	"math"

	"github.com/oakmound/oak/v3/audio/pcm"
)

// Wave functions take a set of options and return an audio
type Wave func(opts ...Option) pcm.Reader

// Sourced from https://en.wikibooks.org/wiki/Sound_Synthesis_Theory/Oscillators_and_Wavetables
func phase(freq Pitch, i int, sampleRate uint32) float64 {
	return float64(freq) * (float64(i) / float64(sampleRate)) * 2 * math.Pi
}

// Sin produces a Sin wave
//         __
//       --  --
//      /      \
//--__--        --__--
func (s Source) Sin(opts ...Option) pcm.Reader {
	return s.Wave(func(s Source, idx int) float64 {
		return s.Volume * math.Sin(s.modPhase(idx))
	}, opts...)
}

// Pulse acts like Square when given a pulse of 2, when given any lesser
// pulse the time up and down will change so that 1/pulse time the wave will
// be up.
//
//     __    __
//     ||    ||
// ____||____||____
func (s Source) Pulse(pulse float64) func(opts ...Option) pcm.Reader {
	pulseSwitch := 1 - 2/pulse
	return func(opts ...Option) pcm.Reader {
		return s.Wave(func(s Source, idx int) float64 {
			if math.Sin(s.Phase(idx)) > pulseSwitch {
				return s.Volume
			}
			return -s.Volume
		}, opts...)
	}
}

// Saw produces a saw wave
//
//   ^   ^   ^
//  / | / | /
// /  |/  |/
func (s Source) Saw(opts ...Option) pcm.Reader {
	return s.Wave(func(s Source, idx int) float64 {
		return s.Volume - (s.Volume / math.Pi * math.Mod(s.Phase(idx), 2*math.Pi))
	}, opts...)
}

// Triangle produces a Triangle wave
//
//   ^   ^
//  / \ / \
// v   v   v
func (s Source) Triangle(opts ...Option) pcm.Reader {
	return s.Wave(func(s Source, idx int) float64 {
		p := s.modPhase(idx)
		m := p * (2 * s.Volume / math.Pi)
		if math.Sin(p) > 0 {
			return -s.Volume + m
		}
		return 3*s.Volume - m
	}, opts...)
}

func (s Source) modPhase(idx int) float64 {
	return math.Mod(s.Phase(idx), 2*math.Pi)
}

// Could have pulse triangle

type Wave8Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int8
}

func (pr *Wave8Reader) ReadPCM(b []byte) (n int, err error) {
	bytesPerI8 := int(pr.Channels)
	for i := 0; i+bytesPerI8 <= len(b); i += bytesPerI8 {
		i8 := pr.waveFunc(pr.Source, pr.lastIndex)
		pr.lastIndex++
		for c := 0; c < int(pr.Channels); c++ {
			b[i+c] = byte(i8)
		}
		n += bytesPerI8
	}
	return
}

func (s Source) Wave(waveFn func(s Source, idx int) float64, opts ...Option) pcm.Reader {
	switch s.Bits {
	case 8:
		s.Volume *= math.MaxInt8
		return &Wave8Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int8 {
				return int8(waveFn(s, idx))
			},
		}
	case 32:
		s.Volume *= math.MaxInt32
		return &Wave32Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int32 {
				return int32(waveFn(s, idx))
			},
		}
	case 16:
		fallthrough
	default:
		s.Volume *= math.MaxInt16
		return &Wave16Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int16 {
				return int16(waveFn(s, idx))
			},
		}
	}
}

type Wave16Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int16
}

func (pr *Wave16Reader) ReadPCM(b []byte) (n int, err error) {
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
	return
}

type Wave32Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int32
}

func (pr *Wave32Reader) ReadPCM(b []byte) (n int, err error) {
	bytesPerF32 := int(pr.Channels) * 4
	for i := 0; i+bytesPerF32 <= len(b); i += bytesPerF32 {
		i32 := pr.waveFunc(pr.Source, pr.lastIndex)
		pr.lastIndex++
		for c := 0; c < int(pr.Channels); c++ {
			b[i+(4*c)] = byte(i32)
			b[i+(4*c)+1] = byte(i32 >> 8)
			b[i+(4*c)+2] = byte(i32 >> 16)
			b[i+(4*c)+3] = byte(i32 >> 24)
		}
		n += bytesPerF32
	}
	return
}
