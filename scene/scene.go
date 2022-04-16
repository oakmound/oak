package scene

import (
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/oakerr"
)

// A Scene is a set of functions defining what needs to happen when a scene
// starts and ends.
type Scene struct {
	// Start is called when a scene begins, including contextual information like
	// what scene came before this one and a direct reference to clean data structures
	// for event handling and rendering
	Start func(ctx *Context)
	// End is a function returning the next scene and a SceneResult of
	// input settings for the next scene.
	End func() (nextScene string, result *Result)
}

// A Result is a set of options for what should be passed into the next
// scene and how the next scene should be transitioned to.
type Result struct {
	NextSceneInput interface{}
	Transition
}

// GoTo returns an End function that, without any other customization possible,
// will change to the input next scene.
func GoTo(nextScene string) func() (nextScene string, result *Result) {
	return func() (string, *Result) {
		return nextScene, nil
	}
}

// GoToPtr returns an End function that, without any other customization possible,
// will change to the input next scene. It takes a pointer so the scene can
// be changed after this function is called.
func GoToPtr(nextScene *string) func() (nextScene string, result *Result) {
	return func() (string, *Result) {
		if nextScene == nil {
			dlog.Error(oakerr.NilInput{InputName: "nextScene"}.Error())
			return "", nil
		}
		return *nextScene, nil
	}
}
