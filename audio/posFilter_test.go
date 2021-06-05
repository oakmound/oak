package audio

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/audio/support/font"
	"github.com/oakmound/oak/v3/audio/support/synth"
)

func TestPosFilter(t *testing.T) {
	kla, err := synth.Int16.Sin()
	if err != nil {
		t.Fatalf("expected sin wave creation to succeed")
	}
	x, y := new(float64), new(float64)
	a := New(DefaultFont, kla.(Data), x, y)
	x2 := 100.0
	y2 := 100.0
	DefaultFont.Filter(PosFilter(NewEars(&x2, &y2, 100, 300)))
	err = <-a.Play()
	if err != nil {
		t.Fatalf("expected playing sin wave to succeed")
	}
	time.Sleep(a.PlayLength())
	// Assert left ear hears audio
	x2 -= 200
	err = <-a.Play()
	if err != nil {
		t.Fatalf("expected playing sin wave (2) to succeed")
	}
	time.Sleep(a.PlayLength())
	// Assert right ear hears audio
	y2 += 500
	err = <-a.Play()
	if err != nil {
		t.Fatalf("expected playing sin wave (3) to succeed")
	}
	time.Sleep(a.PlayLength())
	// Assert nothing is heard
	*DefaultFont = *font.New()
	DefaultFont.Filter(PosFilter(NewEars(&x2, &y2, 100, 2000)))
	x2 -= 200
	err = <-a.Play()
	if err != nil {
		t.Fatalf("expected playing sin wave (4) to succeed")
	}
	time.Sleep(a.PlayLength())
	// Assert right ear hears audio
	x2 += 1000
	err = <-a.Play()
	if err != nil {
		t.Fatalf("expected playing sin wave (5) to succeed")
	}
	time.Sleep(a.PlayLength())
	// Assert left ear hears audio

	_, _ = kla.Filter(PosFilter(NewEars(&x2, &y2, 0, 0)))
	// assert.NotNil(t, err)
}
