package mouse

// Mouse events: MousePress, MouseRelease, MouseScrollDown, MouseScrollUp, MouseDrag
// Payload: (mouse.Event) details on the mouse event
const (
	Press      = "MousePress"
	Release    = "MouseRelease"
	ScrollDown = "MouseScrollDown"
	ScrollUp   = "MouseScrollUp"
	Click      = "MouseClick"
	Drag       = "MouseDrag"
	//
	PressOn      = Press + "On"
	ReleaseOn    = Release + "On"
	ScrollDownOn = ScrollDown + "On"
	ScrollUpOn   = ScrollUp + "On"
	ClickOn      = Click + "On"
	DragOn       = Drag + "On"
)
