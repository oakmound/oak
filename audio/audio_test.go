package audio

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/audio/internal/audio/filter"
	"github.com/oakmound/oak/v3/audio/internal/synth"
)

func TestAudioFuncs(t *testing.T) {
	kla, err := synth.Int16.Sin()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := New(DefaultFont, kla.(Data))
	err = a.SetVolume(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = <-a.Play()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(a.PlayLength())
	// Assert audio is playing
	<-a.Play()
	err = a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(a.PlayLength())
	// Assert audio is not playing
	kla, err = a.Copy()
	a = kla.(*Audio)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert audio is playing
	a = a.MustCopy().(*Audio)
	if a.Xp() != nil {
		t.Fatalf("audio without position had x pointer")
	}
	if a.Yp() != nil {
		t.Fatalf("audio without position had y pointer")
	}
	kla, err = a.Filter(filter.Volume(.5))
	a = kla.(*Audio)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert quieter audio is playing
	a = a.MustFilter(filter.Volume(.5)).(*Audio)
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert yet quieter audio is playing
	a.SetVolume(-2000)
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert yet quieter audio is playing

}
