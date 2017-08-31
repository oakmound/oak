package oakerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsAreErrors(t *testing.T) {

	var err error
	err = NotLoaded{}
	assert.NotEmpty(t, err.Error())
	err = ExistingElement{}
	assert.NotEmpty(t, err.Error())
	err = ExistingElement{Overwritten: true}
	assert.NotEmpty(t, err.Error())
	err = InsufficientInputs{}
	assert.NotEmpty(t, err.Error())
	// Assert nothing crashed
}
