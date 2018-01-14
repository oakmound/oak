package oakerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsAreErrors(t *testing.T) {

	var err error = NotLoaded{}
	assert.NotEmpty(t, err.Error())
	err = ExistingElement{}
	assert.NotEmpty(t, err.Error())
	err = ExistingElement{Overwritten: true}
	assert.NotEmpty(t, err.Error())
	err = InsufficientInputs{}
	assert.NotEmpty(t, err.Error())
	err = InvalidInput{}
	assert.NotEmpty(t, err.Error())
	err = NilInput{}
	assert.NotEmpty(t, err.Error())
	err = IndivisibleInput{}
	assert.NotEmpty(t, err.Error())
	err = IndivisibleInput{IsList: true}
	assert.NotEmpty(t, err.Error())
	err = ConsError{ExistingElement{}, ExistingElement{}}
	assert.NotEmpty(t, err.Error())
	err = UnsupportedFormat{}
	assert.NotEmpty(t, err.Error())
	err = InvalidLength{}
	assert.NotEmpty(t, err.Error())
	err = UnsupportedPlatform{}
	assert.NotEmpty(t, err.Error())
	// Assert nothing crashed
}
