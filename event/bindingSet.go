package event

// A Mapping stores a slice of event names and bindings
type Mapping struct {
	eventNames []string
	binds      []Bindable
}

// A BindingSet stores sets of event mappings bound to string names.
// The use case for a BindingSet is for a character that can exist in multiple states,
// so that they can swiftly switch between the event bindings that define those
// states.
type BindingSet map[string]Mapping

// Set makes a new EventMapping for BindingSet
func (b BindingSet) Set(setName string, mappingSets ...map[string]Bindable) BindingSet {

	numMappings := 0
	for _, m := range mappingSets {
		numMappings += len(m)

	}
	bindings := make([]Bindable, numMappings)
	events := make([]string, numMappings)
	i := 0
	for _, m := range mappingSets {
		for k, v := range m {
			bindings[i] = v
			events[i] = k
			i++
		}
	}

	b[setName] = Mapping{eventNames: events, binds: bindings}
	return b
}

// RebindMapping resets the entity controlling this cid to only have the bindings
// in the passed in event mapping
func (cid CID) RebindMapping(mapping Mapping) {
	cid.UnbindAllAndRebind(mapping.binds, mapping.eventNames)
}
