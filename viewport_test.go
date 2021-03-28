package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/scene"
)

func sleep() {
	// TODO V3: test how far we can bring this down and get consistent results
	time.Sleep(300 * time.Millisecond)
}

func TestViewport(t *testing.T) {
	c1 := NewController()
	err := c1.SceneMap.Add("blank",
		// Initialization function
		func(*scene.Context) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) { return "blank", nil })
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	go c1.Init("blank")
	time.Sleep(2 * time.Second)
	if (c1.ViewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{0, 0})
	}
	c1.SetScreen(5, 5)
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{5, 5}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{5, 5})
	}
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 4, 4))
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{5, 5}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{5, 5})
	}
	c1.SetScreen(-1, -1)
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{0, 0})
	}
	c1.SetScreen(6, 6)
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{0, 0})
	}
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 1000, 1000))
	c1.SetScreen(20, 20)
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{20, 20}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{20, 20})
	}
	c1.SetViewportBounds(intgeom.NewRect2(21, 21, 2000, 2000))
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{21, 21}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{21, 21})
	}
	c1.SetScreen(1000, 1000)
	sleep()
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 900, 900))
	sleep()
	if (c1.ViewPos) != (intgeom.Point2{900 - c1.Width(), 900 - c1.Height()}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{900 - c1.Width(), 900 - c1.Height()})
	}

	c1.skipSceneCh <- true

	sleep()

	if (c1.ViewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.ViewPos, intgeom.Point2{0, 0})
	}
}
