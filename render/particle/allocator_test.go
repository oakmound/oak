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

func TestAllocatorLookup(t *testing.T) {
	src := NewSource(NewColorGenerator(), 0)
	cid := src.CID
	pidBlock := Allocate(cid)
	src2 := LookupSource(pidBlock * blockSize)
	assert.Equal(t, src, src2)
	assert.Nil(t, Lookup((pidBlock*blockSize)+1))
	Deallocate(2)
	Deallocate(1)
	Deallocate(0)
}
