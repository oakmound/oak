package mouse

import "github.com/oakmound/oak/v3/event"

// Binding will convert a function that accepts a typecast *mouse.Event into a generic event binding
//
// Example:
// 		bus.Bind(mouse.ClickOn, mouse.Binding(clickHandler))
func Binding(fn func(event.CID, *Event) int) func(event.CID, interface{}) int {
	return func(cid event.CID, iface interface{}) int {
		me, ok := iface.(*Event)
		if !ok {
			return event.UnbindSingle
		}
		return fn(cid, me)
	}
}
