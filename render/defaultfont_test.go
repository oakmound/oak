package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLegacyFont(t *testing.T) {
	initTestFont()
	assert.NotNil(t, NewText(dummyStringer{}, 0, 0))
	assert.NotNil(t, NewStrText("text", 0, 0))
	assert.NotNil(t, NewIntText(new(int), 0, 0))
}
