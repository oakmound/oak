package scene

import (
	"sync"

	"github.com/oakmound/oak/v4/oakerr"
)

// A Map lets scenes be accessed via associated names.
type Map struct {
	CurrentScene string
	scenes       map[string]Scene
	// This could be a RWMutex, but it isn't anticipated that
	// reads will be more common than writes.
	lock sync.Mutex
}

// NewMap creates a scene map
func NewMap() *Map {
	return &Map{
		scenes: map[string]Scene{},
	}
}

// Get returns the scene associated with the given name, if it exists. If it
// does not exist, it returns a zero value and false.
func (m *Map) Get(name string) (Scene, bool) {
	m.lock.Lock()
	s, ok := m.scenes[name]
	m.lock.Unlock()

	return s, ok
}

// GetCurrent returns the current scene, as defined by map.CurrentScene
func (m *Map) GetCurrent() (Scene, bool) {
	return m.Get(m.CurrentScene)
}

// AddScene takes a scene struct, checks that its assigned name does not
// conflict with an existing name in the map, and then adds it to the map.
// If a conflict occurs, the scene will not be overwritten.
// Checks if the Scene's start is nil, sets to noop if so.
// Checks if the Scene's end is nil, sets to loop to this scene if so.
func (m *Map) AddScene(name string, s Scene) error {

	if s.Start == nil {
		s.Start = func(*Context) {}
	}
	if s.End == nil {
		s.End = GoTo(name)
	}

	var err error
	m.lock.Lock()
	if _, ok := m.scenes[name]; ok {
		err = oakerr.ExistingElement{
			InputName:   name,
			InputType:   "scene",
			Overwritten: false,
		}
	} else {
		m.scenes[name] = s
	}
	m.lock.Unlock()

	return err
}
