//go:build linux

package pcm

import (
	"io"
	"fmt"
	"bytes"
	"sync"

	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
)

func initOS() error {
	return nil
}

// This mutex may be unneeded
var newWriterMutex sync.Mutex

func newWriter(f Format) (Writer, error) {
	newWriterMutex.Lock()
	defer newWriterMutex.Unlock()
	// TODO: 
	// 1. Volume scales with pitch-- lower pitches are quieter -- this does not happen on other OSes, so it's probably
	//    not a result of the audio we're generating-- ???
	// 2. If you play too many things too fast the library crashes. What things are concurrently safe and what things aren't?
	//    Can we fix it here, or must we fix it in a fork / patch?
	// 3. read unix @->/run/user/1000/pulse/native: read: connection reset by peer, followed by nil panic, from creating many writers
	channelOpt := pulse.PlaybackStereo 
	if f.Channels == 2 {
		channelOpt = pulse.PlaybackMono
	} else {
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
	pb, err := client.NewPlayback(pulse.NewReader(ignoreEOFReader{handOver}, pfmt), 
		pulse.PlaybackSink(sink),
		pulse.PlaybackSampleRate(int(f.SampleRate)),
		channelOpt, 
		pulse.PlaybackLatency(0.025),
	)
	if err != nil {
		return nil, fmt.Errorf("pulse writer creation failed: %w", err) 
	}

	return &pulseWriter{
		Format: f,
		handOver: handOver, 
		playBack: pb,
	}, nil 
}

type ignoreEOFReader struct {
	io.Reader 
}

func (m ignoreEOFReader) Read(b []byte) (n int, err error) {
	n, err = m.Reader.Read(b)
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
	handOver *bytes.Buffer
	playBack *pulse.PlaybackStream 
	playing bool
}

func (dsw *pulseWriter) Close() error {
	dsw.Lock()
	defer dsw.Unlock()
	var err error
	if dsw.playing {
		dsw.playBack.Stop()
		dsw.playBack.Close()
		dsw.playing = false
	}
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
		go dsw.playBack.Start()
	}
	if err != nil {
		fmt.Println(dsw.playBack.Error())
	}
	return n, err 
}
