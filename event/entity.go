package event

var (
	highestID CID = 0
	callers       = make([]Entity, 0)
)

type Entity interface {
	Init() CID
}

func NextID(e Entity) CID {
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
