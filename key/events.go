package key

const (
	// Down is sent when a key is pressed. It is sent both as
	// Down, and as Down + the key name.
	Down = "KeyDown"
	// Up is sent when a key is released. It is sent both as
	// Up, and as Up + the key name.
	Up = "KeyUp"
	// Held is sent when a key is held down. It is sent both as
	// Held, and as Held + the key name.
	Held = "KeyHeld"
)
