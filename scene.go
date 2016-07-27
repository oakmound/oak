package plastic

import (
	"errors"
	"fmt"
)

var (
	sceneMap    map[string]*Scene = make(map[string]*Scene)
	activeScene *Scene
)

type Scene struct {
	active bool
	start  SceneStart
	loop   SceneLoop
	end    SceneEnd
}

type SceneEnd func() string
type SceneStart func(prevScene string)
type SceneLoop func() bool

func GetScene(s string) *Scene {
	return sceneMap[s]
}

func AddScene(name string, start SceneStart, loop SceneLoop, end SceneEnd) error {
	fmt.Println("[plastic]-------- Adding", name)
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
