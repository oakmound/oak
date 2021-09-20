package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/scene"
)

func sleep() {
	// TODO: test how far we can bring this down and get consistent results
	time.Sleep(300 * time.Millisecond)
}

func TestViewport(t *testing.T) {
	c1 := NewWindow()
	err := c1.SceneMap.AddScene("blank", scene.Scene{})
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	go c1.Init("blank")
	time.Sleep(2 * time.Second)
	if (c1.viewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{0, 0})
	}
	c1.SetScreen(5, 5)
	if (c1.viewPos) != (intgeom.Point2{5, 5}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{5, 5})
	}
	_, ok := c1.GetViewportBounds()
	if ok {
		t.Fatalf("viewport bounds should not be set on scene start")
	}

	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 4, 4))
	if (c1.viewPos) != (intgeom.Point2{5, 5}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{5, 5})
	}
	c1.SetScreen(-1, -1)
	if (c1.viewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{0, 0})
	}
	c1.SetScreen(6, 6)
	if (c1.viewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{0, 0})
	}
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 1000, 1000))
	c1.SetScreen(20, 20)
	if (c1.viewPos) != (intgeom.Point2{20, 20}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{20, 20})
	}
	c1.ShiftScreen(-1, -1)
	if (c1.viewPos) != (intgeom.Point2{19, 19}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{19, 19})
	}
	c1.SetViewportBounds(intgeom.NewRect2(21, 21, 2000, 2000))
	if (c1.viewPos) != (intgeom.Point2{21, 21}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{21, 21})
	}
	c1.SetScreen(1000, 1000)
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 900, 900))
	bds, ok := c1.GetViewportBounds()
	if !ok {
		t.Fatalf("viewport bounds were not enabled")
	}
	if bds != intgeom.NewRect2(0, 0, 900, 900) {
		t.Fatalf("viewport bounds were not set: expected %v got %v", intgeom.NewRect2(0, 0, 900, 900), bds)
	}
	if (c1.viewPos) != (intgeom.Point2{900 - c1.Width(), 900 - c1.Height()}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{900 - c1.Width(), 900 - c1.Height()})
	}
	c1.RemoveViewportBounds()
	_, ok = c1.GetViewportBounds()
	if ok {
		t.Fatalf("viewport bounds were enabled after clear")
	}
	c1.SetViewportBounds(intgeom.NewRect2(0, 0, 900, 900))

	c1.skipSceneCh <- ""

	sleep()

	if (c1.viewPos) != (intgeom.Point2{0, 0}) {
		t.Fatalf("expected %v got %v", c1.viewPos, intgeom.Point2{0, 0})
	}

	_, ok = c1.GetViewportBounds()
	if ok {
		t.Fatalf("viewport bounds should not be set on scene start")
	}
}
