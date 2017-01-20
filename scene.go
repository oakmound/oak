package oak

import (
	"errors"
	"fmt"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	sceneMap = map[string]*Scene{
		"loading": {
			false,
			func(prevScene string, data interface{}) {
				return
			},
			func() bool {
				select {
				case <-startupLoadComplete:
					dlog.Info("Load Complete")
					return false
				default:
					return true
				}
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
