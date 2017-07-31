package audio

import (
	"testing"
	"time"

	"github.com/200sc/klangsynthese/font"
	"github.com/200sc/klangsynthese/synth"
	"github.com/stretchr/testify/assert"
)

func TestPosFilter(t *testing.T) {
	kla, err := synth.Int16.Sin()
	assert.Nil(t, err)
	x, y := new(float64), new(float64)
	a := New(DefFont, kla.(Data), x, y)
	x2 := 100.0
	y2 := 100.0
	DefFont.Filter(PosFilter(NewEars(&x2, &y2, 100, 300)))
	err = <-a.Play()
	time.Sleep(a.PlayLength())
	// Assert left ear hears audio
	x2 -= 200
	err = <-a.Play()
	time.Sleep(a.PlayLength())
	// Assert right ear hears audio
	y2 += 500
	err = <-a.Play()
	time.Sleep(a.PlayLength())
	// Assert nothing is heard
	*DefFont = *font.New()
	DefFont.Filter(PosFilter(NewEars(&x2, &y2, 100, 2000)))
	x2 -= 200
	err = <-a.Play()
	time.Sleep(a.PlayLength())
	// Assert right ear hears audio
	x2 += 1000
	err = <-a.Play()
	time.Sleep(a.PlayLength())
	// Assert left ear hears audio

	_, err = kla.Filter(PosFilter(NewEars(&x2, &y2, 0, 0)))
	assert.NotNil(t, err)
}
