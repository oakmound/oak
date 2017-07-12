package collision

// A Filter will take a set of collision spaces
// and return the subset that match some requirement
type Filter func([]*Space) []*Space

// FirstLabel returns ths first space that has a label in the input, or nothing
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

// WithLabels will only return spaces with a label in the input
func WithLabels(ls ...Label) Filter {
	return func(sps []*Space) []*Space {
		out := make([]*Space, len(sps))
		i := 0
		for _, s := range sps {
			for _, l := range ls {
				if s.Label == l {
					out[i] = s
					i++
				}
			}
		}
		return out[:i+1]
	}
}

// WithoutLabels will return no spaces with a label in the input
func WithoutLabels(ls ...Label) Filter {
	return func(sps []*Space) []*Space {
		out := make([]*Space, len(sps))
		i := 0
	spaceLoop:
		for _, s := range sps {
			for _, l := range ls {
				if s.Label == l {
					continue spaceLoop
				}
			}
			out[i] = s
			i++
		}
		return out[:i+1]
	}
}
