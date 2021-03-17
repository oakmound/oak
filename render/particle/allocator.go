package particle

import "github.com/oakmound/oak/v2/event"

const (
	blockSize = 2048
)

// TODO: add .Stop?
type Allocator struct {
	particleBlocks map[int]event.CID
	nextOpenCh     chan int
	freeCh         chan int
	allocCh        chan event.CID
	requestCh      chan int
	responseCh     chan event.CID
}

func NewAllocator() *Allocator {
	return &Allocator{
		particleBlocks: make(map[int]event.CID),
		nextOpenCh:     make(chan int),
		freeCh:         make(chan int),
		allocCh:        make(chan event.CID),
		requestCh:      make(chan int),
		responseCh:     make(chan event.CID),
	}
}

func (a *Allocator) Run() {
	lastOpen := 0
	for {
		if _, ok := a.particleBlocks[lastOpen]; !ok {
			select {
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
func (a *Allocator) Allocate(id event.CID) int {
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
	return event.GetEntity(owner).(*Source)
}

// Lookup requests a specific particle in the particle space
func (a *Allocator) Lookup(id int) Particle {
	source := a.LookupSource(id)
	return source.particles[id%blockSize]
}
