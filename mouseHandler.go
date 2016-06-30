package plastic

//mousehandler.go

type MouseEvent struct {
	X, Y   float32
	Button string
}

func getMouseButton(i int32) string {
	s := ""
	switch i {
	case 1:
		s = "LeftMouse"
	case 2:
		s = "MiddleMouse"
	case 3:
		s = "RightMouse"
	case -1:
		s = "ScrollUpMouse"
	case -2:
		s = "ScrollDownMouse"
	default:
		s = ""
	}
	return s
}
