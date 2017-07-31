package audio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlay(t *testing.T) {
	_, err := Load(".", "test.wav")
	assert.Nil(t, err)
	err = Play(DefFont, "test.wav")
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	err = DefPlay("test.wav")
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	// Assert something was played twice
	Unload("test.wav")
	err = Play(DefFont, "test.wav")
	assert.NotNil(t, err)
}
