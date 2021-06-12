package debugstream

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestScopedCommands(t *testing.T) {
	sc := NewScopedCommands()
	if len(sc.commands) != 1 {
		t.Fatalf("scoped commands failed to create with one scope: had %v", len(sc.commands))
	}
	if len(sc.commands[0]) != 3 {
		t.Fatalf("scoped commands failed to create with three commands: had %v", len(sc.commands[0]))
	}
}

func TestScopedCommands_AttachToStream(t *testing.T) {
	in := bytes.NewBufferString("simple")
	out := new(bytes.Buffer)

	sc := NewScopedCommands()
	sc.AttachToStream(context.Background(), in, out)

	// lazy interim approach for the async to complete

	time.Sleep(50 * time.Millisecond)
	output := out.String()
	if !strings.Contains(output, "Unknown command") {
		t.Fatalf("attached Stream doesnt work %s\n", output)
	}

}

func TestScopedCommands_DetachFromStream(t *testing.T) {
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)

	ctx, cancel := context.WithCancel(context.Background())

	sc := NewScopedCommands()
	sc.AttachToStream(ctx, in, out)
	cancel()
	time.Sleep(50 * time.Millisecond)
	output := out.String()
	if !strings.Contains(output, "stopping debugstream") {
		t.Fatalf("unattaching Stream doesnt work %s\n", output)
	}

}
