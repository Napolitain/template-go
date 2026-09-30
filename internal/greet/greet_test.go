package greet_test

import (
	"testing"

	"github.com/napolitain/template-go/internal/greet"
)

func TestHello(t *testing.T) {
	t.Parallel()

	if got, want := greet.Hello("world"), "Hello, world!"; got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
