package collision

// ReactiveSpace is a space that keeps track of a map of collision events
type ReactiveSpace struct {
	*Space
	onHits map[Label]OnHit
}

// NewEmptyReactiveSpace returns a reactive space with no onHit mapping
func NewEmptyReactiveSpace(s *Space) *ReactiveSpace {
	return &ReactiveSpace{
		Space:  s,
		onHits: make(map[Label]OnHit),
	}
}

// NewReactiveSpace creates a reactive space
func NewReactiveSpace(s *Space, onHits map[Label]OnHit) *ReactiveSpace {
	return &ReactiveSpace{
		Space:  s,
		onHits: onHits,
	}
}

// CallOnHits calls CallOnHits on the underlying space of a reactive space
// with the reactive spaces' map of collision events, and returns the channel
// it will send the done signal from.
func (rs *ReactiveSpace) CallOnHits() chan bool {
	doneCh := make(chan bool)
	go CallOnHits(rs.Space, rs.onHits, doneCh)
	return doneCh
}

// Add adds a mapping to a reactive spaces' onhit map
func (rs *ReactiveSpace) Add(i Label, oh OnHit) {
	rs.onHits[i] = oh
}

// Remove removes a mapping from a reactive spaces' onhit map
func (rs *ReactiveSpace) Remove(i Label) {
	delete(rs.onHits, i)
}

// Clear resets a reactive space's onhit map
func (rs *ReactiveSpace) Clear() {
	rs.onHits = make(map[Label]OnHit)
}
