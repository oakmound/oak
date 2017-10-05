package ray

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
)

// A CastLimit is a function that can be applied to
// a Caster's points to return after it adds each one.
// If a Limit returns false, that Caster will immediately
// cease casting.
//
// If a Caster's ray collides with multiple spaces at the same
// point, and some of them would pass a CastLimit, but others would
// not, a Caster will not reliably return those that would pass the
// limit.
type CastLimit func([]collision.Point) bool

// AddLimit is a helper for converting a CastLimit into a CastOption.
func AddLimit(cl CastLimit) CastOption {
	return func(c *Caster) {
		c.Limits = append(c.Limits, cl)
	}
}

// LimitResults will cause a Caster to return a limited number of
// collision points.
func LimitResults(limit int) CastOption {
	return AddLimit(func(ps []collision.Point) bool {
		return len(ps) < limit
	})
}

// StopAtLabel will cause a caster to cease casting as soon as it
// hits one of the input labels.
func StopAtLabel(ls ...collision.Label) CastOption {
	return AddLimit(func(ps []collision.Point) bool {
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

// StopAtID will cause a caster to cease casting as soon as it
// hits one of the input CIDs.
func StopAtID(ids ...event.CID) CastOption {
	return AddLimit(func(ps []collision.Point) bool {
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
