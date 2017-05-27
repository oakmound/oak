// Package mouse handles the propagation of mouse events
// though clickable regions.
// It has a lot of functions which are equivalent to those in the collision package.
package mouse

import (
	"sync"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"github.com/Sythe2o0/rtreego"
)

var (
	mt *rtreego.Rtree
	// LastMouseEvent is the last triggered mouse event,
	// tracked for continuous mouse responsiveness on events
	// that don't take in a mouse event
	LastMouseEvent Event
	addLock        = sync.Mutex{}
)

// An Event is passed in through all Mouse related event bindings to
// indicate what type of mouse event was triggered, where it was triggered,
// and which mouse button it concerns.
// this is a candidate for merging with physics.Vector
type Event struct {
	X, Y   float32
	Button string
	Event  string
}

// ToSpace converts a mouse event into a collision space
func (e Event) ToSpace() *collision.Space {
	return collision.NewUnassignedSpace(float64(e.X), float64(e.Y), 0.1, 0.1)
}

// Init initializes the mouse package with an rtree.
func Init() {
	dlog.Verb("Mouse init started")
	addLock.Lock()
	mt = rtreego.NewTree(20, 40)
	addLock.Unlock()
	dlog.Verb("Mouse init done")
}

// Clear just calls Init
func Clear() {
	dlog.Verb("Mouse clear started ")
	Init()
	dlog.Verb("Mouse clear done")
}

// Add adds a collision space to the mouse rtree
func Add(sp *collision.Space) {
	if sp == nil {
		return
	}
	addLock.Lock()
	mt.Insert(sp)
	addLock.Unlock()
}

// Remove removes a collision space from the mouse rtree.
// Potentially in the future these rtrees and the collision
// rtrees should not be package global items
func Remove(sp *collision.Space) {
	if sp == nil {
		return
	}
	addLock.Lock()
	mt.Delete(sp)
	addLock.Unlock()
}

// UpdateSpace updates the rectangle behind s inside the mouse rtree
func UpdateSpace(x, y, w, h float64, s *collision.Space) {
	if s == nil {
		return
	}
	loc := collision.NewRect(x, y, w, h)
	addLock.Lock()
	mt.Delete(s)
	s.Location = loc
	mt.Insert(s)
	addLock.Unlock()
}

// Hits returns the set of collision spaces intersected by the input space
func Hits(sp *collision.Space) []*collision.Space {
	results := mt.SearchIntersect(sp.Bounds())
	out := make([]*collision.Space, len(results))
	for index, v := range results {
		out[index] = v.(*collision.Space)
	}
	return out
}

// Propagate triggers direct mouse events on entities which are clicked
func Propagate(eventName string, me Event) {
	mouseLoc := collision.NewUnassignedSpace(float64(me.X), float64(me.Y), 0.01, 0.01)
	hits := mt.SearchIntersect(mouseLoc.Bounds())
	for _, v := range hits {
		sp := v.(*collision.Space)
		sp.CID.Trigger(eventName, me)
	}
}

// GetMouseButton is a utitilty function which translates
// integer values of mouse keys from golang's event/mouse library
// into strings.
func GetMouseButton(i int32) (s string) {
	switch i {
	case 1:
		s = "LeftMouse"
	case 2:
		s = "MiddleMouse"
	case 3:
		s = "RightMouse"
	case -1:
		s = "ScrollUpMouse"
	case -2:
		s = "ScrollDownMouse"
	default:
		s = ""
	}
	return
}
