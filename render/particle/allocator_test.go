package particle

import (
	"testing"

	"github.com/oakmound/oak/event"
	"github.com/stretchr/testify/assert"
)

func TestAllocate(t *testing.T) {
	for i := 0; i < 100; i++ {
		assert.Equal(t, Allocate(event.CID(i)), i)
	}
}

func TestDeallocate(t *testing.T) {
	Allocate(0)
	Deallocate(0)
	assert.Equal(t, Allocate(0), 0)
}
