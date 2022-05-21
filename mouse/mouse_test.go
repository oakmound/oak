package mouse

import (
	"testing"

	"golang.org/x/mobile/event/mouse"
)

func TestEventNameIdentity(t *testing.T) {
	if GetEvent(mouse.DirPress, 0) != Press {
		t.Fatalf("event mismatch for event %v, expected %v", mouse.DirPress, "MousePress")
	}
	if GetEvent(mouse.DirRelease, 0) != Release {
		t.Fatalf("event mismatch for event %v, expected %v", mouse.DirRelease, "MouseRelease")
	}
	if GetEvent(mouse.DirNone, -2) != ScrollDown {
		t.Fatalf("event mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollDown")
	}
	if GetEvent(mouse.DirNone, -1) != ScrollUp {
		t.Fatalf("event mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollUp")
	}
	if GetEvent(mouse.DirNone, 0) != Drag {
		t.Fatalf("event mismatch for event %v, expected %v", mouse.DirNone, "MouseDrag")
	}
}
