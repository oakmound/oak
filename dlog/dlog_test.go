package dlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	assert.Equal(t, GetLogLevel(), ERROR)
}
