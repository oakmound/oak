package mouse

// Mouse events: MousePress, MouseRelease, MouseScrollDown, MouseScrollUp, MouseDrag
// Payload: (mouse.Event) details of the mouse event
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

const (
	ButtonLeft      = "LeftMouse"
	ButtonMiddle    = "MiddleMouse"
	ButtonRight     = "RightMouse"
	ButtonWheelUp   = "ScrollUpMouse"
	ButtonWheelDown = "ScrollDownMouse"
)
