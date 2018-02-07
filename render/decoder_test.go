package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDecoder(t *testing.T) {
	assert.NotNil(t, RegisterDecoder("too long", nil))
	short := "s"
	assert.NotNil(t, RegisterDecoder(short, nil))
	assert.NotNil(t, RegisterDecoder(".png", nil))
	assert.Nil(t, RegisterDecoder(".new", nil))
}
