package sequence

import (
	"time"

	"github.com/oakmound/oak/v3/audio/internal/synth"
)

// A ChordPattern represents the order of pitches and holds
// for each of those pitches over a sequence of (potential)
// chords. Todo: pitchPattern is a subset of this, should
// it even exist?
type ChordPattern struct {
	Pitches [][]synth.Pitch
	Holds   [][]time.Duration
}

// HasChords lets generators be built from chord Options
// if they have a pointer to a chord pattern
type HasChords interface {
	GetChordPattern() *ChordPattern
}

// GetChordPattern returns a pointer to a generator's chord pattern
func (cp *ChordPattern) GetChordPattern() *ChordPattern {
	return cp
}

// Chords sets the generator's chord pattern
func Chords(cp ChordPattern) Option {
	return func(g Generator) {
		if hcp, ok := g.(HasChords); ok {
			*(hcp.GetChordPattern()) = cp
		}
	}
}
