package synth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/oakmound/oak/v4/audio"
)

func TestMain(m *testing.M) {
	err := audio.InitDefault()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestFilters(t *testing.T) {
	src := Int16
	// Todo: really gotta fix the sample rate evenness thing
	src.SampleRate = 40000
	src.Volume = .07

	fadeInFrames := time.Second

	unison := 4

	for i := 0; i < unison; i++ {
		go audio.Play(context.Background(), audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw())))
		go audio.Play(context.Background(), audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw(Detune(.04)))))
		go audio.Play(context.Background(), audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw(Detune(-.05)))))
	}
	go audio.Play(context.Background(), audio.FadeIn(fadeInFrames, audio.LoopReader(src.Noise())))

	time.Sleep(3 * time.Second)
}
