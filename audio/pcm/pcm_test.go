package pcm_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/audio/synth"
	"github.com/oakmound/oak/v3/audio/wav"
	"github.com/oakmound/oak/v3/audio/pcm"
)

func TestMain(m *testing.M) {
	err := pcm.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestLoopingWav(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "test.wav"))
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()
	kfmt, err := wav.ReadFormat(f)
	if err != nil {
		t.Fatalf("failed to read wav header in file: %v", err)
	}
	format := pcm.Format{
		SampleRate: kfmt.SampleRate,
		Channels:   kfmt.Channels,
		Bits:       kfmt.Bits,
	}
	w, err := pcm.NewWriter(format)
	if err != nil {
		t.Fatalf("failed to create pcm writer: %v", err)
	}
	r := pcm.LoopReader(&pcm.IOReader{
		Format: format,
		Reader: f,
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = pcm.Play(ctx, w, r)
		if err != nil {
			t.Errorf("failed to play: %v", err)
		}
	}()
	time.Sleep(10 * time.Second)
	fmt.Println("stopping")
	cancel()
	time.Sleep(1 * time.Second)
}

func TestLoopingSin(t *testing.T) {
	format := pcm.Format{
		SampleRate: 44100,
		Channels:   2,
		Bits:       16,
	}
	w, err := pcm.NewWriter(format)
	if err != nil {
		t.Fatalf("failed to create pcm writer: %v", err)
	}

	s := synth.Int16

	s.Volume *= 65535 / 4
	wave := make([]int16, int(s.Seconds*float64(s.SampleRate)))
	for i := 0; i < len(wave); i++ {
		wave[i] = int16(s.Volume * math.Sin(s.Phase(i)))
	}
	b := bytesFromInts(wave, int(s.Channels))
	r := pcm.LoopReader(&pcm.BytesReader{
		Buffer: b,
		Format: format,
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = pcm.Play(ctx, w, r)
		if err != nil {
			t.Errorf("failed to play: %v", err)
		}
	}()
	time.Sleep(10 * time.Second)
	fmt.Println("stopping")
	cancel()
	time.Sleep(1 * time.Second)
}

func bytesFromInts(is []int16, channels int) []byte {
	var ratio = channels * 2
	wave := make([]byte, len(is)*ratio)
	for i := 0; i < len(wave); i += ratio {
		wave[i] = byte(is[i/ratio])
		wave[i+1] = byte(is[i/ratio] >> 8)
		// duplicate the contents across all channels
		for c := 1; c < channels; c++ {
			wave[i+(2*c)] = wave[i]
			wave[i+(2*c)+1] = wave[i+1]
		}
	}
	return wave
}
