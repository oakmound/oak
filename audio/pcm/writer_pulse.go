//go:build linux || darwin

package pcm

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

// This mutex may be unneeded
var newWriterMutex sync.Mutex

func newWriter(f Format) (Writer, error) {
	newWriterMutex.Lock()
	defer newWriterMutex.Unlock()
	// TODO:
	// 1. Volume scales with pitch-- lower pitches are quieter -- only happens for sin32 and tri32, so probably us
	channelOpt := pulse.PlaybackStereo
	if f.Channels == 1 {
		channelOpt = pulse.PlaybackMono
	} else if f.Channels != 2 {
		return nil, fmt.Errorf("unsupported channel count")
	}

	var pfmt byte
	switch f.Bits {
	case 8:
		pfmt = proto.FormatUint8
	case 16:
		pfmt = proto.FormatInt16LE
	case 32:
		pfmt = proto.FormatInt32LE
	default:
		return nil, fmt.Errorf("unsupported bit count")
	}
	client, err := pulse.NewClient()
	if err != nil {
		return nil, err
	}
	sink, err := client.DefaultSink()
	if err != nil {
		return nil, err
	}
	handOver := bytes.NewBuffer([]byte{})
	er := &eofFReader{Buffer: handOver}
	pb, err := client.NewPlayback(pulse.NewReader(er, pfmt),
		pulse.PlaybackSink(sink),
		pulse.PlaybackSampleRate(int(f.SampleRate)),
		channelOpt,
		pulse.PlaybackLatency(0.025),
	)
	if err != nil {
		return nil, fmt.Errorf("pulse writer creation failed: %w", err)
	}

	return &pulseWriter{
		Format:        f,
		handOver:      er,
		playBack:      pb,
		client:        client,
		startComplete: make(chan struct{}),
	}, nil
}

type eofFReader struct {
	*bytes.Buffer
	eof bool
}

func (m *eofFReader) Read(b []byte) (n int, err error) {
	// if we've closed, make sure nothing more is read
	if m.eof {
		return 0, io.EOF
	}
	n, err = m.Buffer.Read(b)
	if err == io.EOF {
		// This enables the pulse library to continually read from
		// the buffer we are handing data over to in WritePCM
		// TODO: to resolve glitches
		err = nil
	}
	return n, err
}

type pulseWriter struct {
	sync.Mutex
	Format
	handOver      *eofFReader
	playBack      *pulse.PlaybackStream
	client        *pulse.Client
	startComplete chan (struct{})
	playing       bool
}

func (dsw *pulseWriter) Close() error {
	dsw.Lock()
	defer dsw.Unlock()
	var err error
	dsw.handOver.eof = true
	if dsw.playing {
		<-dsw.startComplete
		dsw.playBack.Drain()
		dsw.playBack.Stop()
		dsw.playBack.Close()
		dsw.playing = false
	}
	dsw.handOver.Reset()
	dsw.client.Close()
	return err
}

func (dsw *pulseWriter) Reset() error {
	dsw.Lock()
	defer dsw.Unlock()
	// ???
	return nil
}

func (dsw *pulseWriter) WritePCM(data []byte) (n int, err error) {
	dsw.Lock()
	defer dsw.Unlock()
	n, err = dsw.handOver.Write(data)
	if !dsw.playing {
		dsw.playing = true
		go func() {
			// this function blocks if it does not have enough data yet.
			// we don't want to structure our API around start/stop so far,
			// so just start and wait for it to unblock once we've written
			// enough for the os to be happy. Signal to the rest of the system
			// that this process has completed by closing startComplete after.
			// (without this channel, close may be called while start is being called,
			// causing a panic)
			dsw.playBack.Start()
			close(dsw.startComplete)
		}()
	}
	return n, err
}
