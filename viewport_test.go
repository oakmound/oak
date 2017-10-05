package oak

import (
	"image"
	"testing"
	"time"

	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"
)

func testinit() {
	AddScene("blank",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *SceneResult) { return "blank", nil })
	go Init("blank")
	time.Sleep(2 * time.Second)
	// Assert that nothing went wrong
}

func resetOak() {
	select {
	case <-quitCh:
	default:
	}
	lifecycleInit = false
}

func sleep() {
	time.Sleep(300 * time.Millisecond)
}

func TestViewport(t *testing.T) {
	resetOak()
	testinit()
	assert.Equal(t, ViewVector(), physics.NewVector(0, 0))
	assert.Equal(t, ViewPos, image.Point{0, 0})
	SetScreen(5, 5)
	sleep()
	assert.Equal(t, ViewPos, image.Point{5, 5})
	SetViewportBounds(0, 0, 4, 4)
	sleep()
	assert.Equal(t, ViewPos, image.Point{5, 5})
	SetScreen(-1, -1)
	sleep()
	assert.Equal(t, ViewPos, image.Point{0, 0})
	SetScreen(6, 6)
	sleep()
	assert.Equal(t, ViewPos, image.Point{0, 0})
	SetViewportBounds(0, 0, 1000, 1000)
	SetScreen(20, 20)
	sleep()
	assert.Equal(t, ViewPos, image.Point{20, 20})
	SetViewportBounds(21, 21, 2000, 2000)
	sleep()
	assert.Equal(t, ViewPos, image.Point{21, 21})
	SetScreen(1000, 1000)
	sleep()
	SetViewportBounds(0, 0, 900, 900)
	sleep()
	assert.Equal(t, ViewPos, image.Point{900 - ScreenWidth, 900 - ScreenHeight})

	skipSceneCh <- true

	sleep()

	assert.Equal(t, ViewPos, image.Point{0, 0})
}
