package alg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntRanges(t *testing.T) {
	_, err := NewLinearIntRange(1, 0)

	assert.NotNil(t, err)

	_, err = NewSpreadIntRange(1, -1)

	assert.NotNil(t, err)
}
