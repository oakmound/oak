package audio

import (
	"testing"
	"time"

	"github.com/200sc/klangsynthese/synth"
	"github.com/stretchr/testify/assert"
)

func TestAudioFuncs(t *testing.T) {
	kla, err := synth.Int16.Sin()
	assert.Nil(t, err)
	a := New(DefFont, kla.(Data))
	err = <-a.Play()
	assert.Nil(t, err)
	time.Sleep(a.PlayLength())
	// Assert audio is playing
	<-a.Play()
	err = a.Stop()
	assert.Nil(t, err)
	time.Sleep(a.PlayLength())
	// Assert audio is not playing
	kla, err = a.Copy()
	a = kla.(*Audio)
	assert.Nil(t, err)
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert audio is playing
}

func TestPosFilter(t *testing.T) {
	kla, err := synth.Int16.Sin()
	assert.Nil(t, err)
	x, y := new(float64), new(float64)
	a := New(DefFont, kla.(Data), x, y)
	err = <-a.Play()
	assert.Nil(t, err)
}
