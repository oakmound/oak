package event

// DefaultBus is a global Bus. It uses the DefaultCallerMap internally. It should not be used unless your program is only
// using a single Bus. Preferably multi-bus programs would create their own buses and caller maps specific to each bus's
// use.
var DefaultBus *Bus

// DefaultCallerMap is a global CallerMap. It should not be used unless your program is only using a single CallerMap,
// or in other words definitely only has one event bus running at a time.
var DefaultCallerMap *CallerMap

func init() {
	DefaultCallerMap = NewCallerMap()
	DefaultBus = NewBus(DefaultCallerMap)
}
