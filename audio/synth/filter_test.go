package synth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/audio"
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

	w := audio.MustNewWriter(src.PCMFormat())
	loopSrc := audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw()))
	go audio.Play(context.Background(), w, loopSrc)

	w2 := audio.MustNewWriter(src.PCMFormat())
	loopSrc2 := audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw(Detune(.04))))
	go audio.Play(context.Background(), w2, loopSrc2)

	w3 := audio.MustNewWriter(src.PCMFormat())
	loopSrc3 := audio.FadeIn(fadeInFrames, audio.LoopReader(src.Saw(Detune(.05))))
	go audio.Play(context.Background(), w3, loopSrc3)

	time.Sleep(3 * time.Second)
}
