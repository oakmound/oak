package scene

type Context struct {
	PreviousScene string
	SceneInput    interface{}
	// todo: event bus, collision tree, mouse tree, draw stack, window . . .
}
