package oak

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyStore(t *testing.T) {
	// Todo: maybe not strings, it's debatable what is more useful
	// Strings
	// Pros: Easy to write
	// Cons: Can write wrong thing, "Space" instead of "Spacebar"
	//
	// Enum
	// Pros: Uses less space, probably less time as well
	// Cons: Requires import, key.A instead of "A", keybinds require an extended const block
	setDown("Test")
	assert.True(t, IsDown("Test"))
	setUp("Test")
	assert.False(t, IsDown("Test"))
	go keyHoldLoop()
	setDown("Test")
	time.Sleep(2 * time.Second)
	ok, d := IsHeld("Test")
	assert.True(t, ok)
	assert.True(t, d > 1950*time.Millisecond)
	setUp("Test")
	ok, d = IsHeld("Test")
	assert.False(t, ok)
	assert.True(t, d == 0)

	// KeyBind
	assert.Equal(t, "Test", GetKeyBind("Test"))
	BindKey("Test", "Bound")
	assert.Equal(t, "Bound", GetKeyBind("Test"))
}
