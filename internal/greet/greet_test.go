package greet_test

import (
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/napolitain/template-go/internal/greet"
)

func TestHello(t *testing.T) {
	t.Parallel()

	if got, want := greet.Hello("world"), "Hello, world!"; got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}

// Property-based test: rapid generates names and shrinks failures to a minimal case.
func TestHelloProperties(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(t *rapid.T) {
		name := rapid.String().Draw(t, "name")
		got := greet.Hello(name)

		if !strings.HasPrefix(got, "Hello, ") || !strings.HasSuffix(got, "!") {
			t.Fatalf("Hello(%q) = %q, want \"Hello, <name>!\"", name, got)
		}
		if inner := got[len("Hello, ") : len(got)-1]; inner != name {
			t.Fatalf("Hello(%q) wraps %q, want the name unchanged", name, inner)
		}
	})
}
