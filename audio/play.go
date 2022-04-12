// Package audio provides utilities for playing or writing audio streams to OS consumers
package audio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/oakmound/oak/v3/audio/pcm"
)

// WriterBufferLengthInSeconds defines how much data os-level writers provided by this package will rotate through
// in a theoretical circular buffer.
const WriterBufferLengthInSeconds float64 = .5

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
	FadeOutOnStop time.Duration

	Destination pcm.Writer

	// The span of data that should be copied from reader to writer
	// at a time. If too low, may lose accuracy on windows. If too high,
	// may require manual resets when changing audio sources.
	// Defaults to 50 Milliseconds.
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

	ClearBufferOnStop bool
}

func defaultPlayOptions() PlayOptions {
	return PlayOptions{
		CopyIncrement:     50 * time.Millisecond,
		ChaseIncrements:   2,
		FadeOutOnStop:     75 * time.Millisecond,
		ClearBufferOnStop: true,
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
func Play(ctx context.Context, src pcm.Reader, options ...PlayOption) error {
	opts := defaultPlayOptions()
	for _, o := range options {
		o(&opts)
	}
	if opts.Destination == nil {
		var err error
		opts.Destination, err = NewWriter(src.PCMFormat())
		if err != nil {
			return err
		}
		defer opts.Destination.Close()
	}
	format := opts.Destination.PCMFormat()
	if !opts.AllowMismatchedFormats {
		if srcFormat := src.PCMFormat(); srcFormat != format {
			return ErrMismatchedPCMFormat
		}
	}
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
		_, err = opts.Destination.WritePCM(buf)
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	tick := time.NewTicker(opts.CopyIncrement)
	defer tick.Stop()
	// Once we're done, keep writing empty data until the buffer is cleared, unless told not to
	// Do not clear this buffer immediately! You will clear audio data that is actively playing, which will clip!
	if opts.ClearBufferOnStop {
		defer func() {
			buf = make([]byte, format.BytesPerSecond()/uint32(time.Second/opts.CopyIncrement))
			for totalDur := time.Duration(0); totalDur < time.Duration(float64(time.Second)*WriterBufferLengthInSeconds); totalDur += opts.CopyIncrement {
				<-tick.C
				opts.Destination.WritePCM(buf)
			}
		}()
	}
	for {
		select {
		case <-ctx.Done():
			if opts.FadeOutOnStop == 0 {
				return nil
			} else {
				src = FadeOut(opts.FadeOutOnStop, src)
				stopAt := time.NewTimer(opts.FadeOutOnStop * 2)
				defer stopAt.Stop()
				for {
					select {
					case <-stopAt.C:
						return nil
					case <-tick.C:
						_, err := ReadFull(src, buf)
						if errors.Is(err, io.EOF) {
							return nil
						}
						if err != nil {
							return fmt.Errorf("failed to read: %w", err)
						}
						_, err = opts.Destination.WritePCM(buf)
						if err != nil {
							return fmt.Errorf("failed to write: %w", err)
						}
					}
				}
			}
		case <-tick.C:
			_, err := ReadFull(src, buf)
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read: %w", err)
			}
			_, err = opts.Destination.WritePCM(buf)
			if err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
		}
	}
}
