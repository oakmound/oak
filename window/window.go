package window

type Window interface {
	SetFullScreen(bool) error
	SetBorderless(bool) error
	SetTopMost(bool) error
	SetTitle(bool) error
	SetTrayIcon(string) error
	ShowNotification()
	MoveWindow(x, y, w, h int32) error
	GetMonitorSize() (int, int)
	Close() error
}
