package event

var DefaultBus *Bus

// DefaultCallerMap is the caller map used by all event package caller
// functions.
var DefaultCallerMap *CallerMap

func init() {
	DefaultCallerMap = NewCallerMap()
	DefaultBus = NewBus(DefaultCallerMap)
}
