package scene_test

import (
	"fmt"

	"github.com/oakmound/oak/v2/scene"
)

func ExampleMap_GetCurrent() {
	m := scene.NewMap()
	sc := scene.Scene{
		Start: func(prevScene string, data interface{}) {
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
		sc.Start("screen0", nil)
	}
	// Output: Screen two did not exist
	// Starting screen one
}
