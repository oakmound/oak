package timing

import (
	"context"
	"testing"
	"time"
)

func TestDoAfterCancels(t *testing.T) {
	triggered := false
	go DoAfter(3*time.Second, func() {
		triggered = true
	})
	// Wait to make sure the routine started
	time.Sleep(1 * time.Second)
outer:
	for {
		select {
		case ClearDelayCh <- true:
		default:
			break outer
		}
	}
	time.Sleep(3 * time.Second)
	if triggered {
		t.Fatal("doAfter triggered")
	}
}

func TestDoAfterHappens(t *testing.T) {
	triggered := false
	go DoAfter(1*time.Second, func() {
		triggered = true
	})
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatal("doAfter did not trigger")
	}
}

func TestDoAfterContextCancels(t *testing.T) {
	triggered := false
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go DoAfterContext(ctx, func() {
		triggered = true
	})
	// Wait to make sure the routine started
	time.Sleep(1 * time.Second)
outer:
	for {
		select {
		case ClearDelayCh <- true:
		default:
			break outer
		}
	}
	time.Sleep(3 * time.Second)
	if triggered {
		t.Fatal("doAfterContext triggered")
	}
}

func TestDoAfterContextHappens(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	triggered := false
	go DoAfterContext(ctx, func() {
		triggered = true
	})
	time.Sleep(2 * time.Second)
	if !triggered {
		t.Fatal("doAfterContext did not trigger")
	}
}
