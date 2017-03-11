package event

import (
	"fmt"
	"time"

	"bitbucket.org/oakmoundstudio/oak/timing"
)

// Trigger an event, but only
// for one ID. Use case example:
// on onHit event
func (id CID) Trigger(eventName string, data interface{}) chan bool {

	ch := make(chan bool)
	go func(ch chan bool) {

		eb := GetEventBus()
		mutex.RLock()
		if idMap, ok := eb.bindingMap[eventName]; ok {
			if bs, ok := idMap[int(id)]; ok {
				for i := bs.highIndex - 1; i >= 0; i-- {
					for j, bnd := range (*bs.highPriority[i]).sl {
						handleBindable(bnd, int(id), data, j, eventName)
					}
				}
				triggerDefault((bs.defaultPriority).sl, int(id), eventName, data)

				for i := 0; i < bs.lowIndex; i++ {
					for j, bnd := range (*bs.lowPriority[i]).sl {
						handleBindable(bnd, int(id), data, j, eventName)
					}
				}
			}
		}
		mutex.RUnlock()
		select {
		case ch <- true:
		default:
		}
	}(ch)

	return ch
}

func (id CID) TriggerAfter(d time.Duration, eventName string, data interface{}) {
	go func() {
		timing.DoAfter(d, func() {
			id.Trigger(eventName, data)
		})
	}()
}

func Trigger(eventName string, data interface{}) chan bool {
	return thisBus.Trigger(eventName, data)
}

// Called externally by game logic
// and internally by oak itself
// at specific integral points
func (eb_p *EventBus) Trigger(eventName string, data interface{}) chan bool {

	ch := make(chan bool)
	go func(ch chan bool, eb_p *EventBus) {
		eb := (*eb_p)

		mutex.RLock()
		// Loop through all bindableStores for this eventName
		for id, bs := range eb.bindingMap[eventName] {
			// Top to bottom, high priority
			for i := bs.highIndex - 1; i >= 0; i-- {
				for j, bnd := range (*bs.highPriority[i]).sl {
					handleBindable(bnd, id, data, j, eventName)
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
				for j, bnd := range (*bs.lowPriority[i]).sl {
					handleBindable(bnd, id, data, j, eventName)
				}
			}
		}
		mutex.RUnlock()
		select {
		case ch <- true:
		default:
		}
	}(ch, eb_p)

	return ch
}

func triggerDefault(sl []Bindable, id int, eventName string, data interface{}) {
	progCh := make(chan bool)
	for i, bnd := range sl {
		go func(bnd Bindable, id int, eventName string, data interface{}, progCh chan bool, index int) {
			handleBindable(bnd, id, data, index, eventName)
			progCh <- true
		}(bnd, id, eventName, data, progCh, i)
	}
	for range sl {
		<-progCh
	}
}

func handleBindable(bnd Bindable, id int, data interface{}, index int, eventName string) {
	if bnd != nil {
		if id == 0 || GetEntity(id) != nil {
			response := bnd(id, data)
			switch response {
			case UNBIND_EVENT:
				UnbindAll(BindingOption{
					Event{
						eventName,
						id,
					},
					0,
				})
			case UNBIND_SINGLE:
				fmt.Println("Unbinding??")
				Binding{
					BindingOption{
						Event{
							eventName,
							id,
						},
						0,
					},
					index,
				}.Unbind()
			}
		}
	}
}
