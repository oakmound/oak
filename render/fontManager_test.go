package render

import (
	"testing"

	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/oakerr"
	"github.com/stretchr/testify/assert"
)

func TestFontManager(t *testing.T) {

	DefFontGenerator = FontGenerator{}
	fileutil.BindataDir = nil
	fileutil.BindataFn = nil

	fm := NewFontManager()

	f := fm.Get("def")
	assert.NotNil(t, f)
	f = fm.Get("other")
	assert.Nil(t, f)

	err := fm.NewFont("other", FontGenerator{})
	assert.Nil(t, err)

	f = fm.Get("other")
	assert.NotNil(t, f)

	err = fm.NewFont("def", FontGenerator{})
	assert.NotNil(t, err)
	if exists, ok := err.(oakerr.ExistingElement); ok {
		assert.True(t, exists.Overwritten)
	}
}
