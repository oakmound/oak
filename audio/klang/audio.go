// Package klang provides audio playing and encoding support
package klang

import (
	"time"

	"github.com/oakmound/oak/v3/audio/klang/filter/supports"
)

// Audio represents playable, filterable audio data.
type Audio interface {
	// Play returns a channel that will signal when it finishes playing.
	// Looping audio will never send on this channel!
	// The value sent will always be true.
	Play() <-chan error
	// Filter will return an audio with some desired filters applied
	Filter(...Filter) (Audio, error)
	MustFilter(...Filter) Audio
	// Stop will stop an ongoing audio
	Stop() error

	// Implementing struct-- encoding
	Copy() (Audio, error)
	MustCopy() Audio
	PlayLength() time.Duration

	// SetVolume sets the volume of an audio at an OS level,
	// post filters. It multiplies with any volume filters.
	// It takes a value from 0 to -10000, and can only reduce
	// volume from the raw input.
	SetVolume(int32) error
}

// FullAudio supports all the built in filters
type FullAudio interface {
	Audio
	supports.Encoding
	supports.Loop
}

// Stream represents an audio stream. unlike Audio, the length of the
// stream is unknown. Copy is also not supported.
type Stream interface {
	// Play returns a channel that will signal when it finishes playing.
	// Looping audio will never send on this channel!
	// The value sent will always be true.
	Play() <-chan error
	// Filter will return an audio with some desired filters applied
	Filter(...Filter) (Audio, error)
	MustFilter(...Filter) Audio
	// Stop will stop an ongoing audio
	Stop() error
}
