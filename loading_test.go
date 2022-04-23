package oak

import (
	"os"
	"testing"

	"github.com/oakmound/oak/v4/scene"
)

func TestBatchLoad_HappyPath(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.Quit()
		},
	})
	c1.Init("1", func(c Config) (Config, error) {
		c.BatchLoad = true
		c.Assets.AudioPath = "testdata/audio"
		c.Assets.ImagePath = "testdata/images"
		return c, nil
	})
}

func TestBatchLoad_NotFound(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.Quit()
		},
	})
	c1.Init("1", func(c Config) (Config, error) {
		c.BatchLoad = true
		return c, nil
	})
}

func TestBatchLoad_Blank(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.Quit()
		},
	})
	c1.Init("1", func(c Config) (Config, error) {
		c.BatchLoad = true
		c.BatchLoadOptions.BlankOutAudio = true
		return c, nil
	})
}

func TestSetBinaryPayload(t *testing.T) {
	// coverage test, this utility is effectively tested in the render package
	SetFS(os.DirFS("."))
}
