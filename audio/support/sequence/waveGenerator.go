package sequence

import (
	"time"

	"github.com/oakmound/oak/v3/audio/support/audio"
	"github.com/oakmound/oak/v3/audio/support/synth"
)

// A WaveGenerator composes sets of simple waveforms as a sequence
type WaveGenerator struct {
	ChordPattern
	PitchPattern
	WavePattern
	VolumePattern
	HoldPattern
	Length
	Tick
	Loop
}

// NewWaveGenerator uses optional variadic syntax to enable
// any variant of a generator to be made
func NewWaveGenerator(opts ...Option) *WaveGenerator {
	wg := &WaveGenerator{}
	for _, opt := range opts {
		opt(wg)
	}
	return wg
}

// Generate generates a sequence from this wave generator
func (wg *WaveGenerator) Generate() *Sequence {
	sq := &Sequence{}
	sq.Ticker = time.NewTicker(time.Duration(wg.Tick))
	sq.tickDuration = time.Duration(wg.Tick)
	sq.loop = bool(wg.Loop)
	sq.stopCh = make(chan error)
	if wg.Length == 0 {
		if len(wg.PitchPattern) != 0 {
			wg.Length = Length(len(wg.PitchPattern))
		} else if len(wg.ChordPattern.Pitches) != 0 {
			wg.Length = Length(len(wg.ChordPattern.Pitches))
		}
		// else whoops, there's no length
	}
	if len(wg.HoldPattern) == 0 {
		wg.HoldPattern = []time.Duration{sq.tickDuration}
	}
	sq.Pattern = make([]*audio.Multi, wg.Length)

	volumeIndex := 0
	waveIndex := 0
	if len(wg.PitchPattern) != 0 {
		pitchIndex := 0
		holdIndex := 0
		for i := range sq.Pattern {
			p := wg.PitchPattern[pitchIndex]
			if p != synth.Rest {
				a, _ := wg.WavePattern[waveIndex](
					synth.AtPitch(p),
					synth.Duration(wg.HoldPattern[holdIndex]),
					synth.Volume(wg.VolumePattern[volumeIndex]),
				)
				sq.Pattern[i] = audio.NewMulti(a)
			} else {
				sq.Pattern[i] = audio.NewMulti()
			}
			pitchIndex = (pitchIndex + 1) % len(wg.PitchPattern)
			volumeIndex = (volumeIndex + 1) % len(wg.VolumePattern)
			waveIndex = (waveIndex + 1) % len(wg.WavePattern)
			holdIndex = (holdIndex + 1) % len(wg.HoldPattern)
		}
	} else if len(wg.ChordPattern.Pitches) != 0 {
		chordIndex := 0
		for i := range sq.Pattern {
			mult := audio.NewMulti()
			for j, p := range wg.ChordPattern.Pitches[chordIndex] {
				a, _ := wg.WavePattern[waveIndex](
					synth.AtPitch(p),
					synth.Duration(wg.ChordPattern.Holds[chordIndex][j]),
					synth.Volume(wg.VolumePattern[volumeIndex]),
				)
				mult.Audios = append(mult.Audios, a)
			}
			sq.Pattern[i] = mult
			waveIndex = (waveIndex + 1) % len(wg.WavePattern)
			volumeIndex = (volumeIndex + 1) % len(wg.VolumePattern)
			chordIndex = (chordIndex + 1) % len(wg.ChordPattern.Pitches)
		}
	}
	return sq
}
