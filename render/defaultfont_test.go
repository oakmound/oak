package render

import (
	"testing"
)

func TestLegacyFont(t *testing.T) {
	initTestFont()
	if NewStringerText(dummyStringer{}, 0, 0) == nil {
		t.Fatalf("NewStringerText failed")
	}
	if NewText("text", 0, 0) == nil {
		t.Fatalf("NewText failed")
	}
	if NewIntText(new(int), 0, 0) == nil {
		t.Fatalf("NewIntText failed")
	}
	if NewStrPtrText(new(string), 0, 0) == nil {
		t.Fatalf("NewStrPtrText failed")
	}
}
