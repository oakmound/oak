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
	if r == nil {
		return nil, false
	}
	return r.(Renderable), ok
}

// EnumerateDebugRenderableKeys which does not check to see if the associated renderables are still extant
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
