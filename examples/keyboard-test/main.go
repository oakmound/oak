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

type stringStringer string

func (ss stringStringer) String() string {
	return string(ss)
}

func main() {
	oak.Add("keyboard-test", func(string, interface{}) {
		kRenderable := render.NewStrText("", 40, 40)
		render.Draw(kRenderable, 0)
		event.GlobalBind(func(_ int, k interface{}) int {
			kValue := k.(string)
			keyLock.Lock()
			keys[kValue] = struct{}{}
			txt := ""
			for k := range keys {
				txt += k + "\n"
			}
			kRenderable.SetText(stringStringer(txt))
			keyLock.Unlock()
			return 0
		}, key.Down)
		event.GlobalBind(func(_ int, k interface{}) int {
			kValue := k.(string)
			keyLock.Lock()
			delete(keys, kValue)
			txt := ""
			for k := range keys {
				txt += k + " "
			}
			kRenderable.SetText(stringStringer(txt))
			keyLock.Unlock()
			return 0
		}, key.Up)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "keyboard-test", nil
	})
	oak.Init("keyboard-test")
}