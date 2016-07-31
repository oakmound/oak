package collision

type ReactiveSpace struct {
	s      *Space
	onHits map[int]onHit
}

func NewEmptyReactiveSpace(s *Space) *ReactiveSpace {
	return &ReactiveSpace{
		s:      s,
		onHits: make(map[int]onHit),
	}
}

func NewReactiveSpace(s *Space, onHits map[int]onHit) *ReactiveSpace {
	return &ReactiveSpace{
		s:      s,
		onHits: onHits,
	}
}

func (rs *ReactiveSpace) CallOnHits() chan bool {
	doneCh := make(chan bool)
	go CallOnHits(rs.s, rs.onHits, doneCh)
	return doneCh
}

func (rs *ReactiveSpace) Add(i int, oh onHit) {
	rs.onHits[i] = oh
}

func (rs *ReactiveSpace) Remove(i int) {
	delete(rs.onHits, i)
}

func (rs *ReactiveSpace) Clear() {
	rs.onHits = make(map[int]onHit)
}

func (rs *ReactiveSpace) Space() *Space {
	return rs.s
}
