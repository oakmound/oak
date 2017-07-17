package collision

import "github.com/oakmound/oak/event"

// A Filter will take a set of collision spaces
// and return the subset that match some requirement
type Filter func([]*Space) []*Space

// FirstLabel returns the first space that has a label in the input, or nothing
func FirstLabel(ls ...Label) Filter {
	return func(sps []*Space) []*Space {
		for _, s := range sps {
			for _, l := range ls {
				if s.Label == l {
					return []*Space{s}
				}
			}
		}
		return []*Space{}
	}
}

// With will filter spaces so that only those returning true
// from the input keepFn will be in the output
func With(keepFn func(*Space) bool) Filter {
	return func(sps []*Space) []*Space {
		out := make([]*Space, len(sps))
		i := 0
		for _, s := range sps {
			if keepFn(s) {
				out[i] = s
				i++
			}
		}
		return out[:i+1]
	}
}

// Without will filter spaces so that no spaces returning true
// from the input tossFn will be in the output
func Without(tossFn func(*Space) bool) Filter {
	return With(func(s *Space) bool {
		return !tossFn(s)
	})
}

// WithoutCIDs will return no spaces with a CID in the input
func WithoutCIDs(cids ...event.CID) Filter {
	return With(func(s *Space) bool {
		for _, c := range cids {
			if s.CID == c {
				return false
			}
		}
		return true
	})
}

// WithLabels will only return spaces with a label in the input
func WithLabels(ls ...Label) Filter {
	return With(func(s *Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return true
			}
		}
		return false
	})
}

// WithoutLabels will return no spaces with a label in the input
func WithoutLabels(ls ...Label) Filter {
	return With(func(s *Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return false
			}
		}
		return true
	})
}
