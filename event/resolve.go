package event

import (
	"runtime"

	"github.com/oakmound/oak/dlog"
)

// ResolvePending is a contant loop that tracks slices of bind or unbind calls
// and resolves them individually such that they don't break the bus
// Todo: this should be a function on the event bus itself, and should have a better name
// If you ask "Why does this not use select over channels, share memory by communicating",
// the answer is we tried, and it was cripplingly slow.
func ResolvePending() {
	schedCt := 0
	for {
		if len(unbindAllAndRebinds) > 0 {
			resolveUnbindAllAndRebinds()
		}
		// Specific unbinds
		if len(unbinds) > 0 {
			resolveUnbinds()
		}

		// A full set of unbind settings
		if len(fullUnbinds) > 0 {
			resolveFullUnbinds()
		}

		// A partial set of unbind settings
		if len(partUnbinds) > 0 {
			resolvePartialUnbinds()
		}

		// Bindings
		if len(binds) > 0 {
			resolveBindings()
		}

		// This is a tight loop that can cause a pseudo-deadlock
		// by refusing to release control to the go scheduler.
		// This code prevents this from happening.
		// See https://github.com/golang/go/issues/10958
		schedCt++
		if schedCt > 1000 {
			schedCt = 0
			runtime.Gosched()
		}
	}
}

func resolveUnbindAllAndRebinds() {
	mutex.Lock()
	pendingMutex.Lock()
	for _, ubaarb := range unbindAllAndRebinds {
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
			for k := range thisBus.bindingMap {
				namekeys = append(namekeys, k)
			}
		}

		if unbind.CallerID != 0 {
			for _, k := range namekeys {
				delete(thisBus.bindingMap[k], unbind.CallerID)
			}
		} else {
			for _, k := range namekeys {
				delete(thisBus.bindingMap, k)
			}
		}

		dlog.Verb(thisBus.bindingMap)

		// Bindings
		for i := 0; i < len(orderedBindables); i++ {
			fn := orderedBindables[i]
			opt := orderedBindOptions[i]
			list := thisBus.getBindableList(opt)
			list.storeBindable(fn)
		}
	}
	unbindAllAndRebinds = []UnbindAllOption{}
	pendingMutex.Unlock()
	mutex.Unlock()
}

func resolveUnbinds() {
	mutex.Lock()
	pendingMutex.Lock()
	for _, b := range unbinds {
		thisBus.getBindableList(b.BindingOption).removeBinding(b)
	}
	unbinds = []binding{}
	pendingMutex.Unlock()
	mutex.Unlock()
}

func resolveFullUnbinds() {
	mutex.Lock()
	pendingMutex.Lock()
	for _, opt := range fullUnbinds {
		thisBus.getBindableList(opt.BindingOption).removeBindable(opt.Fn)
	}
	fullUnbinds = []UnbindOption{}
	pendingMutex.Unlock()
	mutex.Unlock()
}

func resolvePartialUnbinds() {
	mutex.Lock()
	pendingMutex.Lock()
	for _, opt := range partUnbinds {
		var namekeys []string

		// If we were given a name,
		// we'll just iterate with that name.
		if opt.Name != "" {
			namekeys = append(namekeys, opt.Name)

			// Otherwise, iterate through all events.
		} else {
			for k := range thisBus.bindingMap {
				namekeys = append(namekeys, k)
			}
		}

		if opt.CallerID != 0 {
			for _, k := range namekeys {
				delete(thisBus.bindingMap[k], opt.CallerID)
			}
		} else {
			for _, k := range namekeys {
				delete(thisBus.bindingMap, k)
			}
		}
	}
	partUnbinds = []BindingOption{}
	pendingMutex.Unlock()
	mutex.Unlock()
	dlog.Verb(thisBus.bindingMap)
}

func resolveBindings() {
	mutex.Lock()
	pendingMutex.Lock()
	for _, bindSet := range binds {
		list := thisBus.getBindableList(bindSet.BindingOption)
		list.storeBindable(bindSet.Fn)
	}
	binds = []UnbindOption{}
	pendingMutex.Unlock()
	mutex.Unlock()
}
