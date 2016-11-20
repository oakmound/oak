package collision

type ReactiveSpace struct {
	S      *Space
	onHits map[int]OnHit
}

func NewEmptyReactiveSpace(s *Space) *ReactiveSpace {
	return &ReactiveSpace{
		S:      s,
		onHits: make(map[int]OnHit),
	}
}

func NewReactiveSpace(s *Space, onHits map[int]OnHit) *ReactiveSpace {
	return &ReactiveSpace{
		S:      s,
		onHits: onHits,
	}
}

func (rs *ReactiveSpace) CallOnHits() chan bool {
	doneCh := make(chan bool)
	go CallOnHits(rs.S, rs.onHits, doneCh)
	return doneCh
}

func (rs *ReactiveSpace) Add(i int, oh OnHit) {
	rs.onHits[i] = oh
}

func (rs *ReactiveSpace) Remove(i int) {
	delete(rs.onHits, i)
}

func (rs *ReactiveSpace) Clear() {
	rs.onHits = make(map[int]OnHit)
}

func (rs *ReactiveSpace) Space() *Space {
	return rs.S
}

func (rs *ReactiveSpace) SetDim(w, h float64) {
	rs.S.Update(rs.S.GetX(), rs.S.GetY(), w, h)
}
