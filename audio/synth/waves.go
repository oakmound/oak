// Package synth provides functions and types to support waveform synthesis.
package synth

import (
	"math"
	"math/rand"

	"github.com/oakmound/oak/v4/audio/pcm"
)

// Wave functions take a set of options and return an audio
type Wave func(opts ...Option) pcm.Reader

// Sourced from https://en.wikibooks.org/wiki/Sound_Synthesis_Theory/Oscillators_and_Wavetables
func phase(freq Pitch, i int, sampleRate uint32) float64 {
	return float64(freq) * (float64(i) / float64(sampleRate)) * 2 * math.Pi
}

// Sin produces a Sin wave
//             __
//           --  --
//          /      \
//    --__--        --__--
func (s Source) Sin(opts ...Option) pcm.Reader {
	return s.Wave(Source.SinWave, opts...)
}

func (s Source) SinWave(idx int) float64 {
	return s.Volume * math.Sin(s.modPhase(idx))
}

func (s Source) Square(opts ...Option) pcm.Reader {
	return s.Pulse(2)(opts...)
}

// Pulse acts like Square when given a pulse of 2, when given any lesser
// pulse the time up and down will change so that 1/pulse time the wave will
// be up.
//
//         __    __
//         ||    ||
//     ____||____||____
func (s Source) Pulse(pulse float64) func(opts ...Option) pcm.Reader {
	return func(opts ...Option) pcm.Reader {
		return s.Wave(PulseWave(pulse), opts...)
	}
}

func PulseWave(pulse float64) Waveform {
	pulseSwitch := 1 - 2/pulse
	return func(s Source, idx int) float64 {
		if math.Sin(s.Phase(idx)) > pulseSwitch {
			return s.Volume
		}
		return -s.Volume
	}
}

// Saw produces a saw wave
//
//       ^   ^   ^
//      / | / | /
//     /  |/  |/
func (s Source) Saw(opts ...Option) pcm.Reader {
	return s.Wave(Source.SawWave, opts...)
}

func (s Source) SawWave(idx int) float64 {
	return s.Volume - (s.Volume / math.Pi * math.Mod(s.Phase(idx), 2*math.Pi))
}

// Triangle produces a Triangle wave
//
//       ^   ^
//      / \ / \
//     v   v   v
func (s Source) Triangle(opts ...Option) pcm.Reader {
	return s.Wave(Source.TriangleWave, opts...)
}

func (s Source) TriangleWave(idx int) float64 {
	p := s.modPhase(idx)
	m := p * (2 * s.Volume / math.Pi)
	if math.Sin(p) > 0 {
		return -s.Volume + m
	}
	return 3*s.Volume - m
}

// Noise produces random audio data.
func (s Source) Noise(opts ...Option) pcm.Reader {
	return s.Wave(Source.NoiseWave, opts...)
}

var _ Waveform = Source.NoiseWave

// NoiseWave returns noise pcm data bounded by this source's volume.
func (s Source) NoiseWave(idx int) float64 {
	return ((rand.Float64() * 2) - 1) * s.Volume
}

func (s Source) modPhase(idx int) float64 {
	return math.Mod(s.Phase(idx), 2*math.Pi)
}

// A Waveform is a function that can report a point of audio data given some source parameters for generating the audio
// and an index of where in the generated waveform the requested point lies
type Waveform func(s Source, idx int) float64

// Wave converts a waveform function into a pcm.Reader
func (s Source) Wave(waveFn Waveform, opts ...Option) pcm.Reader {
	switch s.Bits {
	case 8:
		s.Volume *= math.MaxInt8
		return &wave8Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int8 {
				return int8(waveFn(s, idx))
			},
		}
	case 32:
		s.Volume *= math.MaxInt32
		return &wave32Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int32 {
				return int32(waveFn(s, idx))
			},
		}
	case 16:
		fallthrough
	default:
		s.Volume *= math.MaxInt16
		return &wave16Reader{
			Source: s.Update(opts...),
			waveFunc: func(s Source, idx int) int16 {
				return int16(waveFn(s, idx))
			},
		}
	}
}

// MultiWave converts a series of waveform functions into a combined reader, outputting the average
// of all of the source waveforms at any given index
func (s Source) MultiWave(waveFns []Waveform, opts ...Option) pcm.Reader {
	return s.Wave(func(s Source, idx int) float64 {
		var out float64
		for _, wv := range waveFns {
			v := wv(s, idx)
			out += v / float64(len(waveFns))
		}
		return out
	}, opts...)
}

type wave8Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int8
}

func (pr *wave8Reader) ReadPCM(b []byte) (n int, err error) {
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

type wave16Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int16
}

func (pr *wave16Reader) ReadPCM(b []byte) (n int, err error) {
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

type wave32Reader struct {
	Source
	lastIndex int
	waveFunc  func(s Source, idx int) int32
}

func (pr *wave32Reader) ReadPCM(b []byte) (n int, err error) {
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
