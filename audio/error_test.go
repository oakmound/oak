package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oakmound/oak/oakerr"
)

func TestErrorChannel(t *testing.T) {
	err := oakerr.ExistingFontError{}
	err2 := <-errChannel(err)
	assert.Equal(t, err, err2)
}
