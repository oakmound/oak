package particle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseGenerator(t *testing.T) {
	bg := new(BaseGenerator)
	bg.setDefaults()
	bg.ShiftX(10)
	bg.ShiftY(10)
	assert.Equal(t, bg.X(), 10.0)
	assert.Equal(t, bg.X(), 10.0)
}
