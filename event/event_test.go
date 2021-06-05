package event

import (
	"testing"
	"time"
)

func sleep() {
	// this is effectively "sync", or wait for the previous
	// goroutine job to get its job done (Trigger
	// use channels for 'done' signals because we don't want
	// to enable users to wait on triggers that won't actually
	// happen -because they are waiting- within a call that is
	// holding a lock)
	time.Sleep(200 * time.Millisecond)
}

func TestBus(t *testing.T) {
	triggers := 0
	go ResolveChanges()
	GlobalBind("T", Empty(func() {
		triggers++
	}))
	sleep()
	<-TriggerBack("T", nil)
	if triggers != 1 {
		t.Fatalf("first trigger did not happen")
	}
	Trigger("T", nil)
	sleep()
	if triggers != 2 {
		t.Fatalf("second trigger did not happen")
	}
}

func TestUnbind(t *testing.T) {
	triggers := 0
	go ResolveChanges()
	GlobalBind("T", func(CID, interface{}) int {
		triggers++
		return UnbindSingle
	})
	sleep()
	<-TriggerBack("T", nil)
	if triggers != 1 {
		t.Fatalf("first trigger did not happen")
	}
	sleep()
	Trigger("T", nil)
	sleep()
	if triggers != 1 {
		t.Fatalf("second trigger after unbind happened")
	}
	GlobalBind("T", func(CID, interface{}) int {
		triggers++
		return 0
	})
	GlobalBind("T", func(CID, interface{}) int {
		triggers++
		return UnbindEvent
	})
	sleep()
	Trigger("T", nil)
	sleep()
	if triggers != 3 {
		t.Fatalf("global triggers did not happen")
	}

	Trigger("T", nil)
	sleep()
	if triggers != 3 {
		t.Fatalf("global triggers happened after unbind")
	}

	GlobalBind("T", func(CID, interface{}) int {
		triggers++
		return 0
	})
	sleep()
	Trigger("T", nil)
	sleep()
	if triggers != 4 {
		t.Fatalf("global triggers did not happen")
	}

	Reset()

	Trigger("T", nil)
	sleep()
	if triggers != 4 {
		t.Fatalf("global triggers did not unbind after reset")
	}
}

type ent struct{}

func (e ent) Init() CID {
	return NextID(e)
}

func TestCID(t *testing.T) {
	triggers := 0
	go ResolveChanges()
	cid := CID(0).Parse(ent{})
	cid.Bind("T", func(CID, interface{}) int {
		triggers++
		return 0
	})
	sleep()
	cid.Trigger("T", nil)
	sleep()
	if triggers != 1 {
		t.Fatalf("first trigger did not happen")
	}

	// UnbindAllAndRebind
	cid.UnbindAllAndRebind([]Bindable{
		func(CID, interface{}) int {
			triggers--
			return 0
		},
	}, []string{
		"T",
	})

	sleep()
	cid.Trigger("T", nil)
	sleep()
	if triggers != 0 {
		t.Fatalf("second trigger did not happen")
	}

	// UnbindAll
	cid.UnbindAll()

	sleep()
	cid.Trigger("T", nil)
	sleep()
	if triggers != 0 {
		t.Fatalf("second trigger did not unbind")
	}

	cid.Bind("T", func(CID, interface{}) int {
		panic("Should not have been triggered")
	})

	// ResetEntities, etc
	ResetCallerMap()

	cid.Trigger("T", nil)
	sleep()
}

func TestEntity(t *testing.T) {
	go ResolveChanges()
	e := ent{}
	cid := e.Init()
	cid2 := cid.Parse(e)
	if cid != cid2 {
		t.Fatalf("expected id %v got %v", cid, cid2)
	}
	if _, ok := cid.E().(ent); !ok {
		t.Fatalf("cid entity was not present")
	}
	DestroyEntity(cid)
	if cid.E() != nil {
		t.Fatalf("cid entity was not deleted")
	}
}

var (
	ubTriggers int
)

func TestUnbindBindable(t *testing.T) {
	go ResolveChanges()
	GlobalBind("T", tBinding)
	sleep()
	Trigger("T", nil)
	sleep()
	if ubTriggers != 1 {
		t.Fatalf("first trigger did not happen")
	}
	// Fix this syntax
	UnbindBindable(
		UnbindOption{
			Event: Event{
				Name:     "T",
				CallerID: 0,
			},
			Fn: tBinding,
		},
	)
	sleep()
	Trigger("T", nil)
	sleep()
	if ubTriggers != 1 {
		t.Fatalf("unbind call did not unbind trigger")
	}
}

func tBinding(CID, interface{}) int {
	ubTriggers++
	return 0
}

func TestBindableList(t *testing.T) {
	bl := new(bindableList)
	bl.sl = make([]Bindable, 10)
	bl.removeIndex(11)
	bl.sl[2] = tBinding
	bl.removeBindable(tBinding)
	// Assert nothing panicked
}

func TestUnbindAllAndRebind(t *testing.T) {
	go ResolveChanges()
	UnbindAllAndRebind(
		Event{
			Name:     "T",
			CallerID: 0,
		}, []Bindable{}, 0, []string{})
}

func TestBindingSet(t *testing.T) {
	triggers := 0
	bs := BindingSet{}
	bs.Set("one", map[string]Bindable{
		"T": func(CID, interface{}) int {
			triggers++
			return 0
		},
		"P": func(CID, interface{}) int {
			triggers *= 2
			return 0
		},
	})
	e := ent{}
	cid := e.Init()
	cid.RebindMapping(bs["one"])
	sleep()
	cid.Trigger("T", nil)
	sleep()
	cid.Trigger("P", nil)
	sleep()
	if triggers != 2 {
		t.Fatalf("triggers did not happen")
	}
}
