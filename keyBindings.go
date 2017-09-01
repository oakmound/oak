package oak

import (
	"io"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	keyBinds    = make(KeyBindings)
	keyBindLock = sync.RWMutex{}
)

// KeyBindings map input keys to meaningful names, so code can be built around
// those meaningufl names and users can easily rebind which keys do what.
type KeyBindings map[string]string

// Example binding file (toml syntax)
//
// MoveUp = "W"
// MoveDown = "S"
// MoveLeft = "A"
// MoveRight = "D"
// Fire = "Spacebar"

// LoadKeyBindings converts a reader into a map of keys to meaningful names.
// It expects a simple .toml syntax, of key = "value" pairs per line. The resulting
// KeyBindings will have the keys and values reversed, so `MoveUp = "W"` will
// correspond to kb["W"] = "MoveUp"
func LoadKeyBindings(r io.Reader) (KeyBindings, error) {
	kb := make(KeyBindings)
	_, err := toml.DecodeReader(r, &kb)
	return kb, err
}

// BindKeyBindings loads and binds KeyBindings at once. It maintains existing
// keybindings not a part of the input reader.
func BindKeyBindings(r io.Reader) error {
	kb, err := LoadKeyBindings(r)
	if err != nil {
		return err
	}
	BindKeys(kb)
	return nil
}

// SetKeyBindings removes all existing keybindings and then binds all bindings
// within the input reader.
func SetKeyBindings(r io.Reader) error {
	UnbindAllKeys()
	return BindKeyBindings(r)
}

// BindKey binds a name to be triggered when this
// key is triggered
func BindKey(key string, binding string) {
	keyBinds[key] = binding
}

// BindKeys loops over and binds all pairs in the input KeyBindings
func BindKeys(bindings KeyBindings) {
	for k, v := range bindings {
		BindKey(k, v)
	}
}

// UnbindKey removes the binding for the given key in oak's keybindings.
// Does nothing if the key is not already bound.
func UnbindKey(key string) {
	keyBindLock.Lock()
	delete(keyBinds, key)
	keyBindLock.Unlock()
}

// UnbindAllKeys clears the contents of the oak's keybindings.
func UnbindAllKeys() {
	keyBindLock.Lock()
	keyBinds = map[string]string{}
	keyBindLock.Unlock()
}

// GetKeyBind returns either whatever name has been bound to
// a key or the key if nothing has been bound to it.
// Todo: this should be a var function that starts out as "return key",
// and only becomes this function when a binding is made.
func GetKeyBind(key string) string {
	keyBindLock.RLock()
	defer keyBindLock.RUnlock()
	if v, ok := keyBinds[key]; ok {
		return v
	}
	return key
}
