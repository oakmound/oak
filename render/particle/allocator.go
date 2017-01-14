package particle

import "bitbucket.org/oakmoundstudio/oak/event"

const (
	BLOCK_SIZE = 2048
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
					responseCh <- particleBlocks[pID/BLOCK_SIZE]
					lastOpen--
				case i := <-freeCh:
					lastOpen = freeRecieve(i)
				case nextOpenCh <- lastOpen:
					particleBlocks[lastOpen] = <-allocCh
				}
			}
			select {
			case i := <-freeCh:
				lastOpen = freeRecieve(i)
			default:
			}
			lastOpen++
		}
	}()
}

func freeRecieve(i int) int {
	delete(particleBlocks, i)
	return i - 1
}

func Allocate(id event.CID) int {
	nextOpen := <-nextOpenCh
	allocCh <- id
	return nextOpen
}

func Deallocate(block int) {
	freeCh <- block
}

func LookupSource(id int) *Source {
	requestCh <- id
	owner := <-responseCh
	return event.GetEntity(int(owner)).(*Source)
}

func Lookup(id int) Particle {
	source := LookupSource(id)
	return source.particles[id%BLOCK_SIZE]
}
