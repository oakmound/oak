package render

import (
	"testing"
)

func TestDrawStack(t *testing.T) {
	GlobalDrawStack.PreDraw()
	if len(GlobalDrawStack.as) != 1 {
		t.Fatalf("global draw stack did not have one length initially")
	}
	SetDrawStack(
		NewDynamicHeap(),
		NewStaticHeap(),
	)
	if len(GlobalDrawStack.as) != 2 {
		t.Fatalf("global draw stack did not have two length after reset")
	}
	GlobalDrawStack.Pop()
	GlobalDrawStack.PreDraw()
	if len(GlobalDrawStack.as) != 1 {
		t.Fatalf("global draw stack did not have one length after pop")
	}
}
