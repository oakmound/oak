package collision

import "sync"

// ReactiveSpace is a space that keeps track of a map of collision events
type ReactiveSpace struct {
	*Space
	Tree *Tree

	onHitsLock sync.Mutex
	onHits     map[Label]OnHit
}

// NewReactiveSpace creates a reactive space on the default collision tree
func NewReactiveSpace(s *Space, onHits map[Label]OnHit) *ReactiveSpace {
	return &ReactiveSpace{
		Space:  s,
		Tree:   DefaultTree,
		onHits: onHits,
	}
}

// CallOnHits calls CallOnHits on the underlying space of a reactive space
// with the reactive spaces' map of collision events, and returns the channel
// it will send the done signal from. It is not safe to call concurrently with
// add / remove / clear.
func (rs *ReactiveSpace) CallOnHits() chan bool {
	doneCh := make(chan bool)
	go rs.Tree.CallOnHits(rs.Space, rs.onHits, doneCh)
	return doneCh
}

// Add adds a mapping to a reactive spaces' onhit map
func (rs *ReactiveSpace) Add(i Label, oh OnHit) {
	rs.onHitsLock.Lock()
	defer rs.onHitsLock.Unlock()
	rs.onHits[i] = oh
}

// Remove removes a mapping from a reactive spaces' onhit map
func (rs *ReactiveSpace) Remove(i Label) {
	rs.onHitsLock.Lock()
	defer rs.onHitsLock.Unlock()
	delete(rs.onHits, i)
}

// Clear resets a reactive space's onhit map
func (rs *ReactiveSpace) Clear() {
	rs.onHitsLock.Lock()
	defer rs.onHitsLock.Unlock()
	rs.onHits = make(map[Label]OnHit)
}
