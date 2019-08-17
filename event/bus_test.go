package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusStop(t *testing.T) {
	b := NewBus()
	phase := 0
	wait := make(chan struct{})
	go func() {
		assert.Nil(t, b.Stop())
		require.Equal(t, phase, 1)
		wait <- struct{}{}
	}()
	b.updateCh <- true
	<-b.doneCh
	phase = 1
	b.doneCh <- true
	<-wait
}
