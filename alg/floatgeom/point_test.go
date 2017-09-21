package floatgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	assert.Equal(t, a.GreaterOf(b), Point2{1, 1})
}
