package main

import (
	"sync"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var keyLock sync.Mutex
var keys = map[string]struct{}{}

func keyCodeString(c key.Code) string {
	s := c.String()
	return s[4:]
}

func main() {
	oak.AddScene("keyboard-test", scene.Scene{Start: func(*scene.Context) {
		kRenderable := render.NewText("", 40, 40)
		render.Draw(kRenderable, 0)
		event.GlobalBind(key.Down, func(_ event.CID, k interface{}) int {
			kValue := k.(key.Event)
			keyLock.Lock()
			keys[keyCodeString(kValue.Code)] = struct{}{}
			txt := ""
			for k := range keys {
				txt += string(k) + " "
			}
			kRenderable.SetString(txt)
			keyLock.Unlock()
			return 0
		})
		event.GlobalBind(key.Up, func(_ event.CID, k interface{}) int {
			kValue := k.(key.Event)
			keyLock.Lock()
			delete(keys, keyCodeString(kValue.Code))
			txt := ""
			for k := range keys {
				txt += string(k) + " "
			}
			kRenderable.SetString(txt)
			keyLock.Unlock()
			return 0
		})
	}})
	oak.Init("keyboard-test")
}
