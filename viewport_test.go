package oak

import (
	"image"
	"testing"
	"time"

	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/scene"
)

func testinit(t *testing.T) {
	err := SceneMap.Add("blank",
		// Initialization function
		func(*scene.Context) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) { return "blank", nil })
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	go Init("blank")
	time.Sleep(2 * time.Second)
	// Assert that nothing went wrong
}

func resetOak() {
	select {
	case <-quitCh:
	default:
	}
}

func sleep() {
	// TODO V3: test how far we can bring this down and get consistent results
	time.Sleep(300 * time.Millisecond)
}

func TestViewport(t *testing.T) {
	resetOak()
	testinit(t)
	vv := ViewVector()
	if vv.X() != 0 || vv.Y() != 0 {
		t.Fatalf("expected %v got %v", vv, physics.NewVector(0, 0))
	}
	if (ViewPos) != (image.Point{0, 0}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{0, 0})
	}
	SetScreen(5, 5)
	sleep()
	if (ViewPos) != (image.Point{5, 5}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{5, 5})
	}
	SetViewportBounds(0, 0, 4, 4)
	sleep()
	if (ViewPos) != (image.Point{5, 5}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{5, 5})
	}
	SetScreen(-1, -1)
	sleep()
	if (ViewPos) != (image.Point{0, 0}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{0, 0})
	}
	SetScreen(6, 6)
	sleep()
	if (ViewPos) != (image.Point{0, 0}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{0, 0})
	}
	SetViewportBounds(0, 0, 1000, 1000)
	SetScreen(20, 20)
	sleep()
	if (ViewPos) != (image.Point{20, 20}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{20, 20})
	}
	SetViewportBounds(21, 21, 2000, 2000)
	sleep()
	if (ViewPos) != (image.Point{21, 21}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{21, 21})
	}
	SetScreen(1000, 1000)
	sleep()
	SetViewportBounds(0, 0, 900, 900)
	sleep()
	if (ViewPos) != (image.Point{900 - ScreenWidth, 900 - ScreenHeight}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{900 - ScreenWidth, 900 - ScreenHeight})
	}

	skipSceneCh <- true

	sleep()

	if (ViewPos) != (image.Point{0, 0}) {
		t.Fatalf("expected %v got %v", ViewPos, image.Point{0, 0})
	}
}
