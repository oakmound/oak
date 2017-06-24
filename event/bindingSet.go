package event

//EventMapping stores a slice of event names and bindings
type EventMapping struct {
	eventNames []string
	binds      []Bindable
}

// BindingSet maps sets of bindings so that entitys can switch between sets of predefined EventMappings
type BindingSet map[string]EventMapping

//Set makes a new EventMapping for BindingSet
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

	b[setName] = EventMapping{eventNames: events, binds: bindings}
	return b
}

func (c CID) RebindMapping(mapping EventMapping) {
	c.UnbindAllAndRebind(mapping.binds, mapping.eventNames)
}
