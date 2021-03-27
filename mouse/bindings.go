package mouse

import "github.com/oakmound/oak/v2/event"

func Binding(fn func(event.CID, Event) int) func(event.CID, interface{}) int {
	return func(cid event.CID, iface interface{}) int {
		me, ok := iface.(Event)
		if !ok {
			// TODO: log error?
			return event.UnbindSingle
		}
		return fn(cid, me)
	}
}
