package particle

import "github.com/oakmound/oak/event"

const (
	blockSize = 2048
)

var (
	particleBlocks = make(map[int]event.CID)
	nextOpenCh     = make(chan int)
	freeCh         = make(chan int)
	allocCh        = make(chan event.CID)
	requestCh      = make(chan int)
	responseCh     = make(chan event.CID)
)

func init() {
	go func() {
		lastOpen := 0
		for {
			if _, ok := particleBlocks[lastOpen]; !ok {
				select {
				case pID := <-requestCh:
					responseCh <- particleBlocks[pID/blockSize]
					lastOpen--
				case i := <-freeCh:
					opened := freereceive(i)
					if opened < lastOpen {
						lastOpen = opened
					}
				case nextOpenCh <- lastOpen:
					particleBlocks[lastOpen] = <-allocCh
				}
			}
			select {
			case i := <-freeCh:
				opened := freereceive(i)
				if opened < lastOpen {
					lastOpen = opened
				}
			default:
			}
			lastOpen++
		}
	}()
}

func freereceive(i int) int {
	delete(particleBlocks, i)
	return i - 1
}

// Allocate requests a new block in the particle space for the given cid
func Allocate(id event.CID) int {
	nextOpen := <-nextOpenCh
	allocCh <- id
	return nextOpen
}

// Deallocate requests that the given block be removed from the particle space
func Deallocate(block int) {
	freeCh <- block
}

// LookupSource requests the source that generated a pid
func LookupSource(id int) *Source {
	requestCh <- id
	owner := <-responseCh
	return event.GetEntity(int(owner)).(*Source)
}

// Lookup requests a specific particle in the particle space
func Lookup(id int) Particle {
	source := LookupSource(id)
	return source.particles[id%blockSize]
}
