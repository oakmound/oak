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
}

// Init initializes the DirectionSpace as an entity
func (ds *DirectionSpace) Init() event.CID {
	return event.NextID(ds)
}

// NewDirectionSpace creates a DirectionSpace and initializes it as an entity.
func NewDirectionSpace(s *collision.Space, v physics.ForceVector) *DirectionSpace {
	ds := &DirectionSpace{
		Space:       s,
		ForceVector: v,
	}
	s.CID = ds.Init()
	return ds
}
