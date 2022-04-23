package particle

import (
	"testing"

	"github.com/oakmound/oak/v4/event"
)

func TestAllocate(t *testing.T) {
	a := NewAllocator()
	go a.Run()
	for i := 0; i < 100; i++ {
		if a.Allocate(event.CallerID(i)) != i {
			t.Fatalf("expected allocation of id %d to match id", i)
		}
	}
}

func TestDeallocate(t *testing.T) {
	a := NewAllocator()
	go a.Run()

	a.Allocate(0)
	a.Deallocate(0)

	if a.Allocate(0) != 0 {
		t.Fatalf("expected allocation of id %d to match id", 0)
	}
}

func TestAllocatorLookup(t *testing.T) {
	a := NewAllocator()
	go a.Run()

	src := NewDefaultSource(NewColorGenerator(), 0)
	cid := src.CID()
	pidBlock := a.Allocate(cid)
	src2 := a.LookupSource(pidBlock * blockSize)
	if src != src2 {
		t.Fatalf("Lookup on first block did not obtain allocated source")
	}

	src3 := a.Lookup((pidBlock * blockSize) + 1)
	if src3 != nil {
		t.Fatalf("Lookup on second block did not return nil")
	}
	a.Deallocate(2)
	a.Deallocate(1)
	a.Deallocate(0)
}
