package oakerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsAreErrors(t *testing.T) {

	var err error
	err = NotLoadedError{}
	assert.NotEmpty(t, err.Error())
	err = ExistingFontError{}
	assert.NotEmpty(t, err.Error())
	err = InsufficientInputs{}
	assert.NotEmpty(t, err.Error())
	// Assert nothing crashed
}
