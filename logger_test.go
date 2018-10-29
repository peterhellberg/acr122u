package acr122u

import "testing"

func TestStdoutLogger(t *testing.T) {
	l := StdoutLogger()

	if got, want := l.Flags(), 0; got != want {
		t.Fatalf("l.Flags() = %d, want %d", got, want)
	}

	if got, want := l.Prefix(), ""; got != want {
		t.Fatalf("l.Prefix() = %q, want %q", got, want)
	}
}
