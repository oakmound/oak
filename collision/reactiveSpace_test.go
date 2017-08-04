package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReactiveSpace(t *testing.T) {
	Clear()
	var triggered bool
	rs1 := NewEmptyReactiveSpace(NewUnassignedSpace(0, 0, 10, 10))
	assert.NotNil(t, rs1)
	rs2 := NewReactiveSpace(NewUnassignedSpace(5, 5, 10, 10), map[Label]OnHit{
		Label(1): OnIDs(func(id1, id2 int) {
			triggered = true
		}),
	})
	Add(NewLabeledSpace(6, 6, 1, 1, Label(1)))
	<-rs2.CallOnHits()
	assert.True(t, triggered)
	triggered = false

	rs2.Clear()
	<-rs2.CallOnHits()
	assert.False(t, triggered)

	rs1.Add(Label(1), func(*Space, *Space) {
		triggered = true
	})
	<-rs1.CallOnHits()
	assert.True(t, triggered)

	rs1.Remove(Label(1))
	triggered = false
	<-rs1.CallOnHits()
	assert.False(t, triggered)
}
