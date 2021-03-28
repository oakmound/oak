package mouse

import (
	"testing"

	"golang.org/x/mobile/event/mouse"
)

func TestEventNameIdentity(t *testing.T) {
	if GetEventName(mouse.DirPress, 0) != "MousePress" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirPress, "MousePress")
	}
	if GetEventName(mouse.DirRelease, 0) != "MouseRelease" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirRelease, "MouseRelease")
	}
	if GetEventName(mouse.DirNone, -2) != "MouseScrollDown" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollDown")
	}
	if GetEventName(mouse.DirNone, -1) != "MouseScrollUp" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseScrollUp")
	}
	if GetEventName(mouse.DirNone, 0) != "MouseDrag" {
		t.Fatalf("event name mismatch for event %v, expected %v", mouse.DirNone, "MouseDrag")
	}
}
