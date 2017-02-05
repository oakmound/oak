// Package mouse handles the propagation of mouse events
// though clickable regions.
// It has a lot of functions which are equivalent to those in the collision package.
package mouse

import (
	"bitbucket.org/oakmoundstudio/oak/collision"
	"fmt"
	"github.com/Sythe2o0/rtreego"
	"sync"
)

var (
	mt *rtreego.Rtree
	// We track the last triggered mouse event
	// for continuous click hold responsiveness
	LastMouseEvent MouseEvent
	addLock        = sync.Mutex{}
)

type MouseEvent struct {
	X, Y   float32
	Button string
}

func Init() {
	fmt.Println("Mouse init started")
	addLock.Lock()
	mt = rtreego.NewTree(20, 40)
	addLock.Unlock()
	fmt.Println("Mouse init done")
}

func Clear() {
	fmt.Println("Mouse clear started ")
	Init()
	fmt.Println("Mouse clear done")
}

func Add(sp *collision.Space) {
	addLock.Lock()
	mt.Insert(sp)
	addLock.Unlock()
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

func Hits(sp *collision.Space) []*collision.Space {
	results := mt.SearchIntersect(sp.Bounds())
	out := make([]*collision.Space, len(results))
	for index, v := range results {
		out[index] = v.(*collision.Space)
	}
	return out
}

// Trigger direct mouse events on entities
// which are clicked
func Propagate(eventName string, me MouseEvent) {
	mouseLoc := collision.NewUnassignedSpace(float64(me.X), float64(me.Y), 0.01, 0.01)
	hits := mt.SearchIntersect(mouseLoc.Bounds())
	if eventName == "MousePressOn" {
		fmt.Println(len(hits), " Number of mouse collision spaces found")

	}
	for _, v := range hits {
		sp := v.(*collision.Space)
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
