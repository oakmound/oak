package oak

import "testing"

func TestScreenOpts(t *testing.T) {
	// What these functions do (and error presence) depends on the operating
	// system / build tags, which we can't configure at test time without
	// making a new driver just for this test.
	c1 := blankScene(t)
	c1.SetFullScreen(true)
	c1.SetFullScreen(false)
	c1.MoveWindow(10, 10, 20, 20)
	c1.SetBorderless(true)
	c1.SetBorderless(false)
	c1.SetTopMost(true)
	c1.SetTopMost(false)
	c1.SetTitle("testScreenOpts")
	c1.SetTrayIcon("icon.ico")
	c1.ShowNotification("testnotification", "testmessge", true)
	c1.ShowNotification("testnotification", "testmessge", false)
	c1.HideCursor()
}
