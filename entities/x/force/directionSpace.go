package force

import (
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
)

// A DirectionSpace combines collision and a intended direction collision should imply
type DirectionSpace struct {
	*collision.Space
	physics.ForceVector
	event.CallerID
}

func (ds DirectionSpace) CID() event.CallerID {
	return ds.CallerID
}

// NewDirectionSpace creates a DirectionSpace and initializes it as an entity.
func NewDirectionSpace(s *collision.Space, v physics.ForceVector) *DirectionSpace {
	ds := &DirectionSpace{
		Space:       s,
		ForceVector: v,
	}
	// TODO: not default
	s.CID = event.DefaultCallerMap.Register(ds)
	return ds
}
