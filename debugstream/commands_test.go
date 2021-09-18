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

func TestScopedCommands_AssumeScope(t *testing.T) {
	sc := NewScopedCommands()

	in := bytes.NewBufferString("scope\nscope zero\nscope 2\nscope 0\n0 scope 0\n2 scope 0\n0")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(50 * time.Millisecond)

	expected := `assume scope requires a scopeID
Active Scopes: [0]
assume scope <scopeID> expects a valid int32 scope
you provided "zero" which errored with strconv.ParseInt: parsing "zero": invalid syntax
inactive scope 2
assumed scope 0
assumed scope 0
unknown scopeID 2
only provided scopeID 0 without command
`

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
	}
}

func TestScopedCommands_Help(t *testing.T) {
	sc := NewScopedCommands()

	in := bytes.NewBufferString("help\nhelp 0\n help scope\nhelp 1\nhelp badcommand")
	out := new(bytes.Buffer)

	sc.AttachToStream(context.Background(), in, out)

	time.Sleep(50 * time.Millisecond)

	expected := `help <scopeID> to see commands linked to a given window
Active Scopes: [0]
Current Assumed Scope: 0
General Commands:
  fade: fade the specified renderable by the given int if given. Renderable must be registered in debugtools
  help
  scope: provide a scopeID to use commands without a scopeID prepended

help <scopeID> to see commands linked to a given window
Active Scopes: [0]
Current Assumed Scope: 0
General Commands:
  fade: fade the specified renderable by the given int if given. Renderable must be registered in debugtools
  help
  scope: provide a scopeID to use commands without a scopeID prepended

help <scopeID> to see commands linked to a given window
Active Scopes: [0]
Current Assumed Scope: 0
Registered Instances of scope
  scope0 scope: provide a scopeID to use commands without a scopeID prepended
inactive scope 1 see correct usage by using help without the scope
help <scopeID> to see commands linked to a given window
Active Scopes: [0]
Current Assumed Scope: 0
Registered Instances of badcommand
  Warning scope '0' did not have the specified command "badcommand"
`

	got := out.String()
	if got != expected {
		t.Fatal("got:\n" + got + "\nexpected:\n" + expected)
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
