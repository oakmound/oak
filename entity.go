package plastic

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
)

var (
	highestID event.CID = 0
	callers             = make([]Entity, 0)
)

type Entity interface {
	Init() event.CID
}

func NextID(e Entity) event.CID {
	highestID++
	callers = append(callers, e)
	return highestID
}

func GetEntity(i int) interface{} {
	return callers[i-1]
}

func DestroyEntity(i int) {
	callers[i-1] = nil
}
