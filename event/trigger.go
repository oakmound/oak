package event

// Trigger an event, but only
// for one ID. Use case example:
// on onHit event
func (id CID) Trigger(eventName string, data interface{}) {
	eb := GetEventBus()

	mutex.RLock()
	rLocks++
	if idMap, ok := eb.bindingMap[eventName]; ok {
		if bs, ok := idMap[int(id)]; ok {
			for i := bs.highIndex - 1; i >= 0; i-- {
				for _, bnd := range (*bs.highPriority[i]).sl {
					if bnd != nil {
						response := bnd(int(id), data)
						switch response {
						case UNBIND_EVENT:
							thisBus.bindingMap[eventName][int(id)].highPriority[i] = (new(BindableList))
						}
					}
				}
			}
			triggerDefault((bs.defaultPriority).sl, int(id), eventName, data)

			for i := 0; i < bs.lowIndex; i++ {
				for _, bnd := range (*bs.lowPriority[i]).sl {
					if bnd != nil {
						response := bnd(int(id), data)
						switch response {
						case UNBIND_EVENT:
							thisBus.bindingMap[eventName][int(id)].lowPriority[i] = (new(BindableList))
						}
					}
				}
			}
		}
	}
	rLocks--
	mutex.RUnlock()
}

// Called externally by game logic
// and internally by plastic itself
// at specific integral points
func (eb_p *EventBus) Trigger(eventName string, data interface{}) {
	eb := (*eb_p)

	mutex.RLock()
	rLocks++
	// Loop through all bindableStores for this eventName
	for id, bs := range eb.bindingMap[eventName] {
		// Loop through all bindables
		if bs == nil {
			continue
		}

		// Top to bottom, high priority
		for i := bs.highIndex - 1; i >= 0; i-- {
			for _, bnd := range (*bs.highPriority[i]).sl {
				if bnd != nil {
					response := bnd(id, data)
					switch response {
					case UNBIND_EVENT:
						thisBus.bindingMap[eventName][id].highPriority[i] = (new(BindableList))
					}
				}
			}
		}
	}

	for id, bs := range eb.bindingMap[eventName] {
		if bs != nil && bs.defaultPriority != nil {
			triggerDefault((bs.defaultPriority).sl, id, eventName, data)
		}
	}

	//mutex.Lock()
	for id, bs := range eb.bindingMap[eventName] {
		// Bottom to top, low priority
		for i := 0; i < bs.lowIndex; i++ {
			for _, bnd := range (*bs.lowPriority[i]).sl {
				if bnd != nil {
					response := bnd(id, data)
					switch response {
					case UNBIND_EVENT:
						thisBus.bindingMap[eventName][id].lowPriority[i] = (new(BindableList))
					}
				}
			}
		}
	}
	rLocks--
	mutex.RUnlock()
}

func Trigger(eventName string, data interface{}) {
	thisBus.Trigger(eventName, data)
}

func triggerDefault(sl []Bindable, id int, eventName string, data interface{}) {
	progCh := make(chan bool)
	for _, bnd := range sl {
		go func(bnd Bindable, id int, eventName string, data interface{}, progCh chan bool) {
			if bnd != nil {
				if id == 0 || GetEntity(id) != nil {
					response := bnd(id, data)
					switch response {
					case UNBIND_EVENT:
						thisBus.UnbindAll(BindingOption{
							Event{
								eventName,
								id,
							},
							0,
						})
					}
				}
			}
			progCh <- true
		}(bnd, id, eventName, data, progCh)
	}
	for range sl {
		<-progCh
	}
}
