package mouse

import (
	"testing"

	"github.com/oakmound/oak/v3/event"
)

func TestEventOn(t *testing.T) {
	t.Run("AllEvents", func(t *testing.T) {
		if ev2, ok := EventOn(Press); !ok || ev2 != PressOn {
			t.Error("Press was not matched to PressOn")
		}
		if ev2, ok := EventOn(Release); !ok || ev2 != ReleaseOn {
			t.Error("Release was not matched to ReleaseOn")
		}
		if ev2, ok := EventOn(ScrollDown); !ok || ev2 != ScrollDownOn {
			t.Error("ScrollDown was not matched to ScrollDownOn")
		}
		if ev2, ok := EventOn(ScrollUp); !ok || ev2 != ScrollUpOn {
			t.Error("ScrollUp was not matched to ScrollUpOn")
		}
		if ev2, ok := EventOn(Click); !ok || ev2 != ClickOn {
			t.Error("Click was not matched to ClickOn")
		}
		if ev2, ok := EventOn(Drag); !ok || ev2 != DragOn {
			t.Error("Drag was not matched to DragOn")
		}
	})
	t.Run("Unknown", func(t *testing.T) {
		ev := event.RegisterEvent[*Event]()
		_, ok := EventOn(ev)
		if ok {
			t.Error("EventOn should have returned false for an unknown event")
		}
	})
}

func TestEventRelative(t *testing.T) {
	t.Run("AllEvents", func(t *testing.T) {
		if ev2, ok := EventRelative(PressOn); !ok || ev2 != RelativePressOn {
			t.Error("PressOn was not matched to RelativePressOn")
		}
		if ev2, ok := EventRelative(ReleaseOn); !ok || ev2 != RelativeReleaseOn {
			t.Error("ReleaseOn was not matched to RelativeReleaseOn")
		}
		if ev2, ok := EventRelative(ScrollDownOn); !ok || ev2 != RelativeScrollDownOn {
			t.Error("ScrollDownOn was not matched to RelativeScrollDownOn")
		}
		if ev2, ok := EventRelative(ScrollUpOn); !ok || ev2 != RelativeScrollUpOn {
			t.Error("ScrollUpOn was not matched to RelativeScrollUpOn")
		}
		if ev2, ok := EventRelative(ClickOn); !ok || ev2 != RelativeClickOn {
			t.Error("ClickOn was not matched to RelativeClickOn")
		}
		if ev2, ok := EventRelative(DragOn); !ok || ev2 != RelativeDragOn {
			t.Error("DragOn was not matched to RelativeDragOn")
		}
	})
	t.Run("Unknown", func(t *testing.T) {
		ev := event.RegisterEvent[*Event]()
		_, ok := EventRelative(ev)
		if ok {
			t.Error("EventRelative should have returned false for an unknown event")
		}
	})
}
