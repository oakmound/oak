// Package sequence provides generators and options for creating audio sequences.
package sequence

import (
	"time"

	"github.com/oakmound/oak/v3/audio/pcm"
)

// A Sequence is a timed pattern of simultaneously played audios.
type Sequence struct {
	// Sequences play patterns of audio
	// everything at Pattern[0] will be simultaneously Play()ed at
	// Sequence.Play()
	Pattern      []pcm.Reader
	patternIndex int
	// Every tick, the next index in Pattern will be played by a Sequence
	// until the pattern is over.
	Ticker *time.Ticker
	// needed to copy Ticker
	// consider: replacing ticker with dynamic ticker
	tickDuration time.Duration
	stopCh       chan error
	loop         bool
}
