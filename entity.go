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
	if HasEntity(i) {
		return callers[i-1]
	}
	return nil
}

func HasEntity(i int) bool {
	return len(callers) >= i
}

func DestroyEntity(i int) {
	callers[i-1] = nil
}

func ResetEntities() {
	highestID = 0
	callers = []Entity{}
}
