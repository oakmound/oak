package collision

import "github.com/oakmound/oak/event"

type CastLimit func([]Point) bool

func LimitResults(limit int) CastOption {
	return AddLimit(func(ps []Point) bool {
		return len(ps) < limit
	})
}

func StopAtLabel(ls ...Label) CastOption {
	return AddLimit(func(ps []Point) bool {
		z := ps[len(ps)-1].Zone
		if z == nil {
			return true
		}
		for _, l := range ls {
			if z.Label == l {
				return false
			}
		}
		return true
	})
}

func StopAtID(ids ...event.CID) CastOption {
	return AddLimit(func(ps []Point) bool {
		z := ps[len(ps)-1].Zone
		if z == nil {
			return true
		}
		for _, id := range ids {
			if z.CID == id {
				return false
			}
		}
		return true
	})
}

func AddLimit(cl CastLimit) CastOption {
	return func(c *Caster2) {
		c.Limits = append(c.Limits, cl)
	}
}
