package collision

// An OnHit is a function which takes in two spaces
type OnHit func(s, s2 *Space)

// CallOnHits will send a signal to the passed in channel
// when it has completed all collision functions in the hitmap.
func CallOnHits(s *Space, onHits map[Label]OnHit, doneCh chan bool) {
	DefaultTree.CallOnHits(s, onHits, doneCh)
}

// CallOnHits will send a signal to the passed in channel
// when it has completed all collision functions in the hitmap.
func (t *Tree) CallOnHits(s *Space, onHits map[Label]OnHit, doneCh chan bool) {
	progCh := make(chan bool)
	hits := t.Hits(s)
	for _, s2 := range hits {
		go func(s, s2 *Space, onHits map[Label]OnHit, progCh chan bool) {
			if fn, ok := onHits[s2.Label]; ok {
				fn(s, s2)
				progCh <- true
				return
			}
			progCh <- false
		}(s, s2, onHits, progCh)
	}
	// This waits to send our signal that we've
	// finished until we've counted signals for
	// each collision entity
	hitFlag := false
	for range hits {
		v := <-progCh
		hitFlag = hitFlag || v
	}
	doneCh <- hitFlag
}

// OnIDs converts a function on two CIDs to an OnHit
func OnIDs(fn func(int, int)) func(s, s2 *Space) {
	return func(s, s2 *Space) {
		fn(int(s.CID), int(s2.CID))
	}
}
