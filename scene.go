package oak

import (
	"errors"
	"fmt"
)

var (
	sceneMap = map[string]*Scene{
		"loading": {
			false,
			func(prevScene string, data interface{}) {
				return
			},
			func() bool {
				return startupLoadComplete
			},
			func() (string, interface{}) {
				return globalFirstScene, nil
			},
		},
	}
	activeScene *Scene
)

type Scene struct {
	active bool
	start  SceneStart
	loop   SceneLoop
	end    SceneEnd
}

type SceneEnd func() (string, interface{})
type SceneStart func(prevScene string, data interface{})
type SceneLoop func() bool

func GetScene(s string) *Scene {
	return sceneMap[s]
}

func AddScene(name string, start SceneStart, loop SceneLoop, end SceneEnd) error {
	fmt.Println("[oak]-------- Adding", name)
	if _, ok := sceneMap[name]; !ok {
		sceneMap[name] = &(Scene{
			false,
			start,
			loop,
			end,
		})
		return nil
	}
	return errors.New("The scene " + name + " is already defined.")
}
