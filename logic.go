package oak

import "bitbucket.org/oakmoundstudio/oak/timing"

var (
	// LogicTicker is exposed so that games can manually change the speed
	// at which EnterFrame events are produced
	LogicTicker *timing.DynamicTicker

	framesElapsed int
)

// FramesElapsed returns the number of logical frames
// that have been processed in the current scene. This
// value is also passed in to all EnterFrame bindings.
func FramesElapsed() int {
	return framesElapsed
}

func logicLoop() chan bool {
	// The logical loop.
	// In order, it waits on receiving a signal to begin a logical frame.
	// It then runs any functions bound to when a frame begins.
	// It then allows a scene to perform it's loop operation.
	ch := make(chan bool)
	framesElapsed = 0
	go func(doneCh chan bool) {
		LogicTicker = timing.NewDynamicTicker()
		LogicTicker.SetTick(timing.FPSToDuration(FrameRate))
		for {
			select {
			case <-LogicTicker.C:
				<-eb.TriggerBack("EnterFrame", framesElapsed)
				framesElapsed++
				sceneCh <- true
			case <-doneCh:
				LogicTicker.Stop()
				return
			}
		}
	}(ch)
	return ch
}
