package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v4/scene"
)

func TestSceneLoopUnknownScene(t *testing.T) {
	c1 := NewWindow()
	err := c1.SceneMap.AddScene("blank", scene.Scene{})
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	err = c1.Init("bad")
	if err == nil {
		t.Fatal("expected error from Init on unknown scene")
	}
}

func TestSceneLoopUnknownErrorScene(t *testing.T) {
	c1 := NewWindow()
	err := c1.SceneMap.AddScene("blank", scene.Scene{})
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	c1.ErrorScene = "bad2"
	err = c1.Init("bad")
	if err == nil {
		t.Fatal("expected error from Init to error scene")
	}
}

func TestSceneLoopErrorScene(t *testing.T) {
	c1 := NewWindow()
	err := c1.SceneMap.AddScene("blank", scene.Scene{})
	if err != nil {
		t.Fatalf("Scene Add failed: %v", err)
	}
	c1.ErrorScene = "blank"
	go func() {
		err = c1.Init("bad")
	}()
	time.Sleep(2 * time.Second)
	if err != nil {
		t.Fatalf("error transitioning to unknown scene: %v", err)
	}
}
