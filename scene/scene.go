package scene

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
