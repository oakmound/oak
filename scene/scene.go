package scene

import (
	"github.com/oakmound/oak/v2/dlog"
)

// A Scene is a set of functions defining what needs to happen when a scene
// starts, loops, and ends.
type Scene struct {
	Start
	Loop
	End
}

// A Result is a set of options for what should be passed into the next
// scene and how the next scene should be transitioned to.
type Result struct {
	NextSceneInput interface{}
	Transition
}

// Start is a function taking in a previous scene and a payload
// of data from the previous scene's end.
type Start func(prevScene string, data interface{})

// Loop is a function that returns whether or not the current scene
// should continue to loop.
type Loop func() bool

// End is a function returning the next scene and a SceneResult of
// input settings for the next scene.
type End func() (string, *Result)

// BooleanLoop returns a Loop function that will end a scene as soon as the
// input boolean is false, resetting it to true in the process for the
// next scene
func BooleanLoop(b *bool) Loop {
	return func() bool {
		if !(*b) {
			*b = true
			return false
		}
		return true
	}
}

// GoTo returns an End function that, without any other customization possible,
// will change to the input next scene.
func GoTo(nextScene string) End {
	return func() (string, *Result) {
		return nextScene, nil
	}
}

// GoToPtr returns an End function that, without any other customization possible,
// will change to the input next scene. It takes a pointer so the scene can
// be changed after this function is called.
func GoToPtr(nextScene *string) End {
	return func() (string, *Result) {
		if nextScene == nil {
			dlog.Error("Go To: next scene was nil")
			return "", nil
		}
		return *nextScene, nil
	}
}
