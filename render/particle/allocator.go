package particle

import (
	"github.com/oakmound/oak/v4/event"
)

const (
	blockSize = 2048
)

// An Allocator can allocate ids for particles
type Allocator struct {
	particleBlocks map[int]event.CallerID
	nextOpenCh     chan int
	freeCh         chan int
	allocCh        chan event.CallerID
	requestCh      chan int
	responseCh     chan event.CallerID
	stopCh         chan struct{}
}

// NewAllocator creates a new allocator
func NewAllocator() *Allocator {
	return &Allocator{
		particleBlocks: make(map[int]event.CallerID),
		nextOpenCh:     make(chan int),
		freeCh:         make(chan int),
		allocCh:        make(chan event.CallerID),
		requestCh:      make(chan int),
		responseCh:     make(chan event.CallerID),
		stopCh:         make(chan struct{}),
	}
}

// Run spins up an allocator to accept allocation requests. It will run until
// Stop is called. This is a blocking call.
func (a *Allocator) Run() {
	lastOpen := 0
	for {
		if _, ok := a.particleBlocks[lastOpen]; !ok {
			select {
			case <-a.stopCh:
				return
			case pID := <-a.requestCh:
				a.responseCh <- a.particleBlocks[pID/blockSize]
				lastOpen--
			case i := <-a.freeCh:
				opened := a.freereceive(i)
				if opened < lastOpen {
					lastOpen = opened
				}
			case a.nextOpenCh <- lastOpen:
				a.particleBlocks[lastOpen] = <-a.allocCh
			}
		}
		select {
		case <-a.stopCh:
			return
		case i := <-a.freeCh:
			opened := a.freereceive(i)
			if opened < lastOpen {
				lastOpen = opened
			}
		default:
		}
		lastOpen++
	}
}

// DefaultAllocator is an allocator that starts running as soon as this package is imported.
var DefaultAllocator = NewAllocator()

// This is an always-called init instead of Init because oak does not import this
// package by default. If this package is not used, it will not run this goroutine.
func init() {
	go DefaultAllocator.Run()
}

func (a *Allocator) freereceive(i int) int {
	delete(a.particleBlocks, i)
	return i - 1
}

// Allocate requests a new block in the particle space for the given cid
func (a *Allocator) Allocate(id event.CallerID) int {
	nextOpen := <-a.nextOpenCh
	a.allocCh <- id
	return nextOpen
}

// Deallocate requests that the given block be removed from the particle space
func (a *Allocator) Deallocate(block int) {
	a.freeCh <- block
}

// LookupSource requests the source that generated a pid
func (a *Allocator) LookupSource(id int) *Source {
	a.requestCh <- id
	owner := <-a.responseCh
	// TODO: not default?
	return event.DefaultCallerMap.GetEntity(owner).(*Source)
}

// Lookup requests a specific particle in the particle space
func (a *Allocator) Lookup(id int) Particle {
	source := a.LookupSource(id)
	return source.particles[id%blockSize]
}

// Stop stops the allocator's ongoing Run. Once stopped, allocator may not be reused.
// Stop must not be called more than once.
func (a *Allocator) Stop() {
	close(a.stopCh)
}
