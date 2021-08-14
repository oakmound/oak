package debugtools

import (
	"github.com/oakmound/oak/v3/render"
	"golang.org/x/sync/syncmap"
)

var (
	debugMap syncmap.Map
)

// SetDebugRenderable stores a renderable under a name in a package global map.
// this is used by some built in debugConsole helper functions.
func SetDebugRenderable(rName string, r render.Renderable) {
	debugMap.Store(rName, r)
}

// GetDebugRenderable returns whatever renderable is stored under the input
// string, if any.
func GetDebugRenderable(rName string) (render.Renderable, bool) {
	r, ok := debugMap.Load(rName)
	if r == nil {
		return nil, false
	}
	return r.(render.Renderable), ok
}

// EnumerateDebugRenderableKeys lists all registered renderables by key.
// It does not check to see if the associated renderables are still valid in any respect.
func EnumerateDebugRenderableKeys() []string {
	keys := []string{}
	debugMap.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if ok {
			keys = append(keys, key)
		}
		return true
	})
	return keys
}
