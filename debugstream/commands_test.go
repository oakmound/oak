package debugstream

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestScopedCommands(t *testing.T) {
	sc := NewScopedCommands()
	if len(sc.commands) > 0 {
		t.Fatalf("scoped commands failed to create as expected")
	}

}

func TestScopedCommands_AttachToStream(t *testing.T) {
	in := bytes.NewBufferString("simple")
	out := new(bytes.Buffer)

	AttachToStream(in, out)

	// lazy interim approach for the async to complete

	time.Sleep(50 * time.Millisecond)
	output := out.String()
	if !strings.Contains(output, "Unknown command") {
		t.Fatalf("attached Stream doesnt work %s\n", output)
	}

}

func TestScopedCommands_UnAttachFromStream(t *testing.T) {
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)

	AttachToStream(in, out)
	DefaultCommands.UnAttachFromStream()
	time.Sleep(50 * time.Millisecond)
	output := out.String()
	if !strings.Contains(output, "stopping debugstream") {
		t.Fatalf("unattaching Stream doesnt work %s\n", output)
	}

}
