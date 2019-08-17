package audio

import (
	"testing"
	"time"

	"github.com/200sc/go-dist/intrange"
	"github.com/stretchr/testify/assert"
)

func TestChannels(t *testing.T) {
	_, err := DefChannel(intrange.Constant(5))
	assert.NotNil(t, err)
	_, err = Load(".", "test.wav")
	assert.Nil(t, err)
	ch, err := DefChannel(intrange.NewLinear(1, 100), "test.wav")
	assert.Nil(t, err)
	assert.NotNil(t, ch)
	go func() {
		tm := time.Now().Add(2 * time.Second)
		// This only matters when running a suite of tests
		for time.Now().Before(tm) {
			ch <- Signal{0}
		}
	}()
	time.Sleep(2 * time.Second)
}
