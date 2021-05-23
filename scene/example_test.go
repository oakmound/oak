package scene_test

import (
	"fmt"

	"github.com/oakmound/oak/v3/scene"
)

func ExampleMap_GetCurrent() {
	m := scene.NewMap()
	sc := scene.Scene{
		Start: func(*scene.Context) {
			fmt.Println("Starting screen one")
		},
	}
	m.AddScene("screen1", sc)
	m.CurrentScene = "screen2"
	_, ok := m.GetCurrent()
	if !ok {
		fmt.Println("Screen two did not exist")
	}
	m.CurrentScene = "screen1"
	sc, ok = m.GetCurrent()
	if !ok {
		fmt.Println("Screen one did not exist")
	} else {
		sc.Start(&scene.Context{
			PreviousScene: "scene0",
		})
	}
	// Output: Screen two did not exist
	// Starting screen one
}
