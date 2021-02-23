package event

import "time"

// ResolvePending is a constant loop that tracks slices of bind or unbind calls
// and resolves them individually such that they don't break the bus.
// Each section of the loop waits for the predetermined refreshrate prior to attempting to flush.
// Todo: this should have a better name
// If you ask "Why does this not use select over channels, share memory by communicating",
// the answer is we tried, and it was cripplingly slow.
func (eb *Bus) ResolvePending() {
	eb.init.Do(func() {
		go func() {
			for {
				time.Sleep(eb.refreshRate)
				eb.Flush()
			}
		}()
	})
}

// SetRefreshRate on the event bus detailing the time to wait per attempt to ResolvePending.
func (eb *Bus) SetRefreshRate(refreshRate time.Duration) {
	eb.refreshRate = refreshRate
}

func (eb *Bus) resolveUnbindAllAndRebinds() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	for _, ubaarb := range eb.unbindAllAndRebinds {
		unbind := ubaarb.ub
		orderedBindables := ubaarb.bnds
		orderedBindOptions := ubaarb.bs

		var namekeys []string
		// If we were given a name,
		// we'll just iterate with that name.
		if unbind.Name != "" {
			namekeys = append(namekeys, unbind.Name)
			// Otherwise, iterate through all events.
		} else {
			for k := range eb.bindingMap {
				namekeys = append(namekeys, k)
			}
		}

		if unbind.CallerID != 0 {
			for _, k := range namekeys {
				delete(eb.bindingMap[k], unbind.CallerID)
			}
		} else {
			for _, k := range namekeys {
				delete(eb.bindingMap, k)
			}
		}

		// Bindings
		for i := 0; i < len(orderedBindables); i++ {
			fn := orderedBindables[i]
			opt := orderedBindOptions[i]
			list := eb.getBindableList(opt)
			list.storeBindable(fn)
		}
	}
	eb.unbindAllAndRebinds = []UnbindAllOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

func (eb *Bus) resolveUnbinds() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	for _, bnd := range eb.unbinds {
		eb.getBindableList(bnd.BindingOption).removeBinding(bnd)
	}
	eb.unbinds = []binding{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

func (eb *Bus) resolveFullUnbinds() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	for _, opt := range eb.fullUnbinds {
		eb.getBindableList(opt.BindingOption).removeBindable(opt.Fn)
	}
	eb.fullUnbinds = []UnbindOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

func (eb *Bus) resolvePartialUnbinds() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	for _, opt := range eb.partUnbinds {
		var namekeys []string

		// If we were given a name,
		// we'll just iterate with that name.
		if opt.Name != "" {
			namekeys = append(namekeys, opt.Name)

			// Otherwise, iterate through all events.
		} else {
			for k := range eb.bindingMap {
				namekeys = append(namekeys, k)
			}
		}

		if opt.CallerID != 0 {
			for _, k := range namekeys {
				delete(eb.bindingMap[k], opt.CallerID)
			}
		} else {
			for _, k := range namekeys {
				delete(eb.bindingMap, k)
			}
		}
	}
	eb.partUnbinds = []BindingOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}

func (eb *Bus) resolveBindings() {
	eb.mutex.Lock()
	eb.pendingMutex.Lock()
	for _, bindSet := range eb.binds {
		list := eb.getBindableList(bindSet.BindingOption)
		list.storeBindable(bindSet.Fn)
	}
	eb.binds = []UnbindOption{}
	eb.pendingMutex.Unlock()
	eb.mutex.Unlock()
}
