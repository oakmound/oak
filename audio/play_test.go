package audio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayAndLoad(t *testing.T) {
	_, err := Load(".", "test.wav")
	assert.Nil(t, err)
	_, err = Load(".", "badfile.wav")
	assert.NotNil(t, err)
	_, err = Load(".", "play_test.go")
	assert.NotNil(t, err)
	err = Play(DefFont, "test.wav")
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	err = DefPlay("test.wav")
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	// Assert something was played twice
	_, err = GetSounds("test.wav")
	assert.Nil(t, err)
	_, err = GetSounds("badfile.wav")
	assert.NotNil(t, err)
	Unload("test.wav")
	err = Play(DefFont, "test.wav")
	assert.NotNil(t, err)
	err = BatchLoad(".")
	assert.Nil(t, err)
	err = BatchLoad("GarbagePath")
	assert.NotNil(t, err)
}
