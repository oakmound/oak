package oak

import (
	"sync"
	"testing"
	"time"

	"github.com/oakmound/oak/scene"
)

var once sync.Once

func BenchmarkDrawLoop(b *testing.B) {
	once.Do(func() {
		SetupConfig.Debug = Debug{
			"VERBOSE",
			"",
		}
		Add("draw",
			// Initialization function
			func(prevScene string, inData interface{}) {},
			// Loop to continue or stop current scene
			func() bool { return true },
			// Exit to transition to next scene
			func() (nextScene string, result *scene.Result) {
				return "draw", nil
			})
		go Init("draw")
		// give the engine some time to start
		time.Sleep(5 * time.Second)
		// We don't want any regular ticks getting through
		DrawTicker.SetTick(100 * time.Hour)
	})

	b.ResetTimer()
	// This sees how fast the draw ticker will accept forced steps,
	// which won't be accepted until the draw loop itself pulls
	// from the draw ticker, which it only does after having drawn
	// the screen for a frame. This way we push the draw loop
	// to draw as fast as possible and measure that speed.
	for i := 0; i < b.N; i++ {
		DrawTicker.ForceStep()
	}
}
