package collision

import "github.com/oakmound/oak/event"

type CastFilter func(*Space) bool

func AddFilter(cf CastFilter) CastOption {
	return func(c *Caster2) {
		c.Filters = append(c.Filters, cf)
	}
}

func AcceptLabels(ls ...Label) CastOption {
	return AddFilter(func(s *Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return true
			}
		}
		return false
	})
}

func IgnoreLabels(ls ...Label) CastOption {
	return AddFilter(func(s *Space) bool {
		for _, l := range ls {
			if s.Label == l {
				return false
			}
		}
		return true
	})
}

func AcceptIDs(ids ...event.CID) CastOption {
	return AddFilter(func(s *Space) bool {
		for _, id := range ids {
			if s.CID == id {
				return true
			}
		}
		return false
	})
}

func IgnoreIDs(ids ...event.CID) CastOption {
	return AddFilter(func(s *Space) bool {
		for _, id := range ids {
			if s.CID == id {
				return false
			}
		}
		return true
	})
}

func Pierce(n int) CastOption {

	pierced := 0

	return AddFilter(func(s *Space) bool {
		if pierced < n {
			pierced++
			return false
		}
		return true
	})
}
