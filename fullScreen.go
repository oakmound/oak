package oak

type FullScreenable interface {
	SetFullScreen()
}

func SetFullScreen() {
	if fs, ok := windowControl.(FullScreenable); ok {
		fs.SetFullScreen()
	}
}
