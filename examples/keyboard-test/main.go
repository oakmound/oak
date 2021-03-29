package main

import (
	"sync"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

var keyLock sync.Mutex
var keys = map[string]struct{}{}

func main() {
	oak.AddScene("keyboard-test", scene.Scene{Start: func(*scene.Context) {
		kRenderable := render.NewStrText("", 40, 40)
		render.Draw(kRenderable, 0)
		event.GlobalBind(key.Down, func(_ event.CID, k interface{}) int {
			kValue := k.(string)
			keyLock.Lock()
			keys[kValue] = struct{}{}
			txt := ""
			for k := range keys {
				txt += k + "\n"
			}
			kRenderable.SetString(txt)
			keyLock.Unlock()
			return 0
		})
		event.GlobalBind(key.Up, func(_ event.CID, k interface{}) int {
			kValue := k.(string)
			keyLock.Lock()
			delete(keys, kValue)
			txt := ""
			for k := range keys {
				txt += k + " "
			}
			kRenderable.SetString(txt)
			keyLock.Unlock()
			return 0
		})
	}})
	oak.Init("keyboard-test")
}
