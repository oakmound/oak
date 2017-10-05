package ray

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
)

// A CastFilter is a function that can be applied to a Caster
// for each space the Caster's rays hit, returning whether or
// not those spaces should be contained in the Caster's output.
type CastFilter func(*collision.Space) bool

// AddFilter is a utility to convert a CastFilter to a CastOption.
func AddFilter(cf CastFilter) CastOption {
	return func(c *Caster) {
		c.Filters = append(c.Filters, cf)
	}
}

// AcceptLabels signals to a Caster to only return spaces that have
// a space in the set of input labels. If anything in ls is also
// contained by an IgnoreLabels filter on the applied Caster, Ignore
// will dominate.
func AcceptLabels(ls ...collision.Label) CastOption {
	return AddFilter(func(s *collision.Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return true
			}
		}
		return false
	})
}

// IgnoreLabels signals to a Caster to not return spaces that have
// spaces with labels in the set of input labels.
func IgnoreLabels(ls ...collision.Label) CastOption {
	return AddFilter(func(s *collision.Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return false
			}
		}
		return true
	})
}

// AcceptIDs is equivalent to AcceptLabels, but for CIDs.
func AcceptIDs(ids ...event.CID) CastOption {
	return AddFilter(func(s *collision.Space) bool {
		for _, id := range ids {
			if s.CID == id {
				return true
			}
		}
		return false
	})
}

// IgnoreIDs is equivalent to IgnoreLabels, but for CIDs.
func IgnoreIDs(ids ...event.CID) CastOption {
	return AddFilter(func(s *collision.Space) bool {
		for _, id := range ids {
			if s.CID == id {
				return false
			}
		}
		return true
	})
}

// Pierce signals to a Caster to ignore the first n spaces its rays
// collide with, regardless of their composition.
func Pierce(n int) CastOption {

	pierced := 0

	return AddFilter(func(s *collision.Space) bool {
		if pierced < n {
			pierced++
			return false
		}
		return true
	})
}
