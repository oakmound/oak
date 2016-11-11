// Package mouse handles the propagation of mouse events
// though clickable regions.
// It has a lot of functions which are equivalent to those in the collision package.
package mouse

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"github.com/Sythe2o0/rtreego"
)

var (
	mt *rtreego.Rtree
	// We track the last triggered mouse event
	// for continuous click hold responsiveness
	LastMouseEvent MouseEvent
)

type MouseEvent struct {
	X, Y   float32
	Button string
}

func Init() {
	mt = rtreego.NewTree(20, 40)
}

func Clear() {
	Init()
}

func Add(sp *collision.Space) {
	mt.Insert(sp)
}

func Remove(sp *collision.Space) {
	mt.Delete(sp)
}

func UpdateSpace(x, y, w, h float64, s *collision.Space) {
	if s == nil {
		return
	}
	loc := collision.NewRect(x, y, w, h)
	mt.Delete(s)
	s.Location = loc
	mt.Insert(s)
}

// Trigger direct mouse events on entities
// which are clicked
func Propagate(eventName string, me MouseEvent) {
	mouseLoc := collision.NewUnassignedSpace(float64(me.X), float64(me.Y), 0.01, 0.01)
	hits := mt.SearchIntersect(mouseLoc.Bounds())
	for _, v := range hits {
		sp := v.(*collision.Space)

		// Todo:
		// Talk about what event should be triggered here
		sp.CID.Trigger(eventName, nil)
	}
}

// GetMouseButton is a utitilty function which translates
// integer values of mouse keys from golang's event/mouse library
// into strings.
func GetMouseButton(i int32) string {
	s := ""
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
	return s
}
