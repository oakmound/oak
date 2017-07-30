package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func sleep() {
	// this is effectively "sync", or wait for the previous
	// goroutine job to get its job done (event does not
	// use channels for 'done' signals because we don't want
	// to enable users to -wait- on triggers that won't actually
	// happen -because they are waiting- within a call that is
	// holding a lock)
	time.Sleep(200 * time.Millisecond)
}

func TestBus(t *testing.T) {
	triggers := 0
	go ResolvePending()
	GlobalBind(func(int, interface{}) int {
		triggers++
		return 0
	}, "T")
	sleep()
	<-TriggerBack("T", nil)
	assert.Equal(t, triggers, 1)
	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 2)
}

func TestUnbind(t *testing.T) {
	triggers := 0
	go ResolvePending()
	GlobalBind(func(int, interface{}) int {
		triggers++
		return UnbindSingle
	}, "T")
	sleep()
	<-TriggerBack("T", nil)
	assert.Equal(t, triggers, 1)
	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 1)

	GlobalBind(func(int, interface{}) int {
		triggers++
		return 0
	}, "T")
	GlobalBind(func(int, interface{}) int {
		triggers++
		return UnbindEvent
	}, "T")
	sleep()
	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 3)

	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 3)

	GlobalBind(func(int, interface{}) int {
		triggers++
		return 0
	}, "T")
	sleep()
	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 4)

	ResetBus()

	Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 4)
}

type ent struct{}

func (e ent) Init() CID {
	return NextID(e)
}

func TestCID(t *testing.T) {
	triggers := 0
	go ResolvePending()
	cid := CID(0).Parse(ent{})
	cid.Bind(func(int, interface{}) int {
		triggers++
		return 0
	}, "T")
	sleep()
	cid.Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 1)

	// Priority
	var first bool
	cid.BindPriority(func(int, interface{}) int {
		first = true
		return 0
	}, "P", 1)
	cid.BindPriority(func(int, interface{}) int {
		if !first {
			panic("Priority -1 was called before priority 1")
		}
		return 0
	}, "P", -1)
	sleep()
	cid.Trigger("P", nil)

	// UnbindAllAndRebind
	cid.UnbindAllAndRebind([]Bindable{
		func(int, interface{}) int {
			triggers--
			return 0
		},
	}, []string{
		"T",
	})

	sleep()
	cid.Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 0)

	// UnbindAll
	cid.UnbindAll()

	sleep()
	cid.Trigger("T", nil)
	sleep()
	assert.Equal(t, triggers, 0)

	cid.Bind(func(int, interface{}) int {
		panic("Should not have been triggered")
	}, "T")

	// ResetEntities, etc
	assert.Equal(t, cid.String(), "1")
	ResetEntities()

	cid.Trigger("T", nil)
	sleep()
}

func TestEntity(t *testing.T) {
	go ResolvePending()
	e := ent{}
	cid := e.Init()
	cid2 := cid.Parse(e)
	assert.Equal(t, cid, cid2)
	assert.NotNil(t, cid.E().(ent))
	DestroyEntity(int(cid))
	assert.Nil(t, cid.E())
}

var (
	ubTriggers int
)

func TestUnbindBindable(t *testing.T) {
	go ResolvePending()
	GlobalBind(tBinding, "T")
	sleep()
	Trigger("T", nil)
	sleep()
	assert.Equal(t, ubTriggers, 1)
	// Fix this syntax
	UnbindBindable(
		UnbindOption{
			BindingOption: BindingOption{
				Event: Event{
					Name:     "T",
					CallerID: 0,
				},
				Priority: 0,
			},
			Fn: tBinding,
		},
	)
	sleep()
	Trigger("T", nil)
	sleep()
	assert.Equal(t, ubTriggers, 1)
}

func tBinding(int, interface{}) int {
	ubTriggers++
	return 0
}

func TestPriority(t *testing.T) {

}

func TestBindingSet(t *testing.T) {
	triggers := 0
	bs := BindingSet{}
	bs.Set("one", map[string]Bindable{
		"T": func(int, interface{}) int {
			triggers++
			return 0
		},
		"P": func(int, interface{}) int {
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
	assert.Equal(t, triggers, 2)
}
