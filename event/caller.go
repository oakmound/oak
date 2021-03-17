package event

// A Caller can bind, unbind and trigger events.
type Caller interface {
	Trigger(string, interface{})
	Bind(string, Bindable)
	UnbindAll()
	UnbindAllAndRebind([]Bindable, []string)
	E() interface{}
	Parse(Entity) CID
}
