package event

// Oak uses the following built in events:
//
// - CollisionStart/Stop: when a PhaseCollision entity starts/stops touching some label.
//   Payload: (collision.Label) the label the entity has started/stopped touching
//
// - MouseCollisionStart/Stop: as above, for mouse collision
//   Payload: (mouse.Event)
//
// - Mouse events: MousePress, MouseRelease, MouseScrollDown, MouseScrollUp, MouseDrag
//   Payload: (mouse.Event) details on the mouse event
//
// - KeyDown, KeyDown$a: when any key is pressed down, when key $a is pressed down.
//   Payload: (string) the key pressed
//
// - KeyUp, KeyUp$a: when any key is released, when key $a is released.
//   Payload: (string) the key released
//
// And the following:
const (
	// Enter : the beginning of every logical frame.
	// Payload: (int) frames passed since this scene started
	Enter = "EnterFrame"
	//
	// Consider moving the below to different packages before 2.0 release
	//
	// AnimationEnd: Triggered on animations CIDs when they loop from the last to the first frame
	// Payload: nil
	AnimationEnd = "AnimationEnd"
	// ViewportUpdate: Triggered when the position fo of the viewport changes
	// Payload: []float64{viewportX, viewportY}
	ViewportUpdate = "ViewportUpdate"
)
