package pcm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

// WriterBufferLengthInSeconds defines how much data os-level writers provided by this package will rotate through
// in a theoretical circular buffer.
const WriterBufferLengthInSeconds = 1

// InitDefault calls Init with the following value by OS:
// windows: DriverDirectSound
// linux,osx: DriverPulse
func InitDefault() error {
	return Init(DriverDefault)
}

// Init initializes the pcm package to create writer objects with a specific audio driver.
func Init(d Driver) error {
	return initOS(d)
}

// A PlayOption sets some value on a PlayOptions struct.
type PlayOption func(*PlayOptions)

// PlayOptions define ways to configure how playback of some audio proceeds
type PlayOptions struct {
	// The span of data that should be copied from reader to writer
	// at a time. If too low, may lose accuracy on windows. If too high,
	// may require manual resets when changing audio sources.
	// Defaults to 125 Milliseconds.
	CopyIncrement time.Duration
	// How many increments should make up the time between our read and write
	// cursors-- i.e. the audio will be playing at X and we will be writing to
	// X + ChaseIncrements * CopyIncrement.
	// This must be at least 2 to avoid the read and write buffers clipping.
	// Defaults to 2.
	ChaseIncrements int
	// If AllowMismatchedFormats is false, Play will error when a reader's PCM format
	// disagrees with a writer's expected PCM format. Defaults to false.
	AllowMismatchedFormats bool
}

func defaultPlayOptions() PlayOptions {
	return PlayOptions{
		CopyIncrement:   125 * time.Millisecond,
		ChaseIncrements: 2,
	}
}

// ErrMismatchedPCMFormat will be returned by operations streaming from Readers to Writers where the PCM formats
// of those Readers and Writers are not equivalent.
var ErrMismatchedPCMFormat = fmt.Errorf("source and destination have differing PCM formats")

// Play will copy data from the provided src to the provided dst until ctx is cancelled. This copy is not constant.
// The copy will occur in two phases: first, an initial population of the writer to give distance between the read
// cursor and write cursor; immediately upon this write, the writer should begin playback. Following this setup, a
// sub-second amount of data will streamed from src to dst after waiting that same duration. These wait times can
// be configured via PlayOptions.
func Play(ctx context.Context, dst Writer, src Reader, options ...PlayOption) error {
	opts := defaultPlayOptions()
	for _, o := range options {
		o(&opts)
	}
	format := dst.PCMFormat()
	if !opts.AllowMismatchedFormats {
		if srcFormat := src.PCMFormat(); srcFormat != format {
			return ErrMismatchedPCMFormat
		}
	}
	defer dst.Reset()
	buf := make([]byte, format.BytesPerSecond()/uint32(time.Second/opts.CopyIncrement))
	for i := 0; i < opts.ChaseIncrements; i++ {
		// TODO: some formats may expect a minimum buffer size (synth waveforms expect a buffer size of
		// at least bits / 8 * channels), and if the sample rate does not evenly divide that expected minimum,
		// this can hang.
		_, err := ReadFull(src, buf)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}
		_, err = dst.WritePCM(buf)
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	tick := time.NewTicker(opts.CopyIncrement)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			_, err := ReadFull(src, buf)
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read: %w", err)
			}
			_, err = dst.WritePCM(buf)
			if err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
		}
	}
}
