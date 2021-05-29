package debugstream

import "testing"

func TestScopedCommands(t *testing.T) {
	sc := NewScopedCommands()
	if len(sc.commands) != 1 {
		t.Fatalf("scoped commands failed to create as expected")
	}

}
