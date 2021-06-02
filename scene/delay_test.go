package scene

import (
	"context"
	"image/color"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/render"
)

func TestDoAfterCancels(t *testing.T) {
	baseCtx, cancel := context.WithCancel(context.Background())
	ctx := &Context{
		Context: baseCtx,
	}
	triggered := false
	go ctx.DoAfter(3*time.Second, func() {
		triggered = true
	})
	// Wait to make sure the routine started
	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
	if triggered {
		t.Fatal("doAfter should not have triggered")
	}
}

func TestDoAfterHappens(t *testing.T) {
	baseCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx := &Context{
		Context: baseCtx,
	}
	triggered := false
	go ctx.DoAfter(1*time.Second, func() {
		triggered = true
	})
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatal("doAfter did not trigger")
	}
}

func TestDoAfterContextCancels(t *testing.T) {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	ctx := &Context{
		Context: baseCtx,
	}
	triggered := false
	cancelCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go ctx.DoAfterContext(cancelCtx, func() {
		triggered = true
	})
	// Wait to make sure the routine started
	time.Sleep(1 * time.Second)
	baseCancel()
	time.Sleep(3 * time.Second)
	if triggered {
		t.Fatal("doAfterContext should not have triggered")
	}
}

func TestDoAfterContextHappens(t *testing.T) {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()
	ctx := &Context{
		Context: baseCtx,
	}
	cancelCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	triggered := false
	go ctx.DoAfterContext(cancelCtx, func() {
		triggered = true
	})
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatal("doAfterContext did not trigger")
	}
}

func TestDrawForTime(t *testing.T) {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()
	ctx := &Context{
		Context:   baseCtx,
		DrawStack: render.GlobalDrawStack,
	}
	err := ctx.DrawForTime(render.NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 0, 4)
	if err == nil {
		t.Fatalf("draw time to invalid layer should fail")
	}

	err = ctx.DrawForTime(render.NewColorBox(5, 5, color.RGBA{255, 255, 255, 255}), 0, 0)
	if err != nil {
		t.Fatalf("draw time should not have failed")
	}
}
