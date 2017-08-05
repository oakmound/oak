package render

import "golang.org/x/sync/syncmap"

var (
	debugMap syncmap.Map
)

// UpdateDebugMap stores a renderable under a name in a package global map.
// this is used by some built in debugConsole helper functions.
func UpdateDebugMap(rName string, r Renderable) {
	debugMap.Store(rName, r)
}

// GetDebugRenderable returns whatever renderable is stored under the input
// string, if any.
func GetDebugRenderable(rName string) (Renderable, bool) {
	r, ok := debugMap.Load(rName)
	return r.(Renderable), ok
}
