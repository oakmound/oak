package oakerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsAreErrors(t *testing.T) {

	var err error
	err = NotLoaded{}
	assert.NotEmpty(t, err.Error())
	err = ExistingFont{}
	assert.NotEmpty(t, err.Error())
	err = InsufficientInputs{}
	assert.NotEmpty(t, err.Error())
	// Assert nothing crashed
}
