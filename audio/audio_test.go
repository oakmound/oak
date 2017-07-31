package audio

import (
	"testing"
	"time"

	"github.com/200sc/klangsynthese/audio/filter"
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
	a = a.MustCopy().(*Audio)
	assert.Nil(t, a.GetX())
	assert.Nil(t, a.GetY())
	kla, err = a.Filter(filter.Volume(.5))
	a = kla.(*Audio)
	assert.Nil(t, err)
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert quieter audio is playing
	a = a.MustFilter(filter.Volume(.5)).(*Audio)
	a.Play()
	time.Sleep(a.PlayLength())
	// Assert yet quieter audio is playing
}