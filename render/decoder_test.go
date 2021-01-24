package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDecoder(t *testing.T) {
	assert.NotNil(t, RegisterDecoder(".png", nil))
	assert.Nil(t, RegisterDecoder(".new", nil))
}
