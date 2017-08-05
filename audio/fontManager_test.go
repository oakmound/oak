package audio

import (
	"testing"

	"github.com/200sc/klangsynthese/font"
	"github.com/stretchr/testify/assert"
)

func TestFontManager(t *testing.T) {
	fm := NewFontManager()
	assert.Nil(t, fm.NewFont("unused", font.New()))
	assert.NotNil(t, fm.NewFont("unused", font.New()))
	assert.Nil(t, fm.Get("notafont"))
	assert.NotNil(t, fm.Get("def"))
}
