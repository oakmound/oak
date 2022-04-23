package oak

import (
	"context"
	"errors"
	"testing"

	"github.com/oakmound/oak/v4/oakerr"
	"github.com/oakmound/oak/v4/scene"
)

func TestSceneTransition(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			go context.Window.NextScene()
		},
		End: func() (nextScene string, result *scene.Result) {
			return "2", &scene.Result{
				Transition: scene.Fade(1, 10),
			}
		},
	})
	c1.AddScene("2", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.Quit()
		},
	})
	c1.Init("1")
}

func TestLoadingSceneClaimed(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene(oakLoadingScene, scene.Scene{})
	err := c1.Init("1")
	var wantErr oakerr.ExistingElement
	if !errors.As(err, &wantErr) {
		t.Fatalf("expected existing element error, got %v", err)
	}
}

func TestSceneGoTo(t *testing.T) {
	c1 := NewWindow()
	var cancel func()
	c1.ParentContext, cancel = context.WithCancel(c1.ParentContext)
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.GoToScene("good")
		},
		End: func() (nextScene string, result *scene.Result) {
			return "bad", &scene.Result{
				Transition: scene.Fade(1, 10),
			}
		},
	})
	c1.AddScene("good", scene.Scene{
		Start: func(ctx *scene.Context) {
			cancel()
		},
	})
	c1.Init("1")
}
