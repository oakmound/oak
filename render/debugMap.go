package render

import "golang.org/x/sync/syncmap"

var (
	debugMap2 syncmap.Map
)

// UpdateDebugMap sets a value within the debugMap
func UpdateDebugMap(rName string, r Renderable) {
	debugMap2.Store(rName, r)
}

// GetDebugRenderable gets a value from the debugMap
func GetDebugRenderable(rName string) (Renderable, bool) {
	r, ok := debugMap2.Load(rName)
	return r.(Renderable), ok
}
