package acr122u

import (
	"testing"
	"time"

	"github.com/ebfe/scard"
)

func TestEstablishContext(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		scardEstablishContext = func() (*scard.Context, error) {
			return nil, scard.ErrInternalError
		}

		if _, err := EstablishContext(); err != scard.ErrInternalError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		scardEstablishContext = func() (*scard.Context, error) {
			return &scard.Context{}, nil
		}

		if _, err := EstablishContext(); err != scard.ErrInvalidHandle {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestNewContext(t *testing.T) {
	t.Run("Error from IsValid", func(t *testing.T) {
		_, err := newContext(&mockContext{
			isValid: func() (bool, error) {
				return false, scard.ErrInvalidHandle
			},
		})

		if err != scard.ErrInvalidHandle {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Error from ListReaders", func(t *testing.T) {
		_, err := newContext(&mockContext{
			listReaders: func() ([]string, error) {
				return nil, scard.ErrUnknownError
			},
		})

		if err != scard.ErrUnknownError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("No Readers Available", func(t *testing.T) {
		_, err := newContext(&mockContext{
			listReaders: func() ([]string, error) {
				return nil, nil
			},
		})

		if err != scard.ErrNoReadersAvailable {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		ctx, err := newContext(&mockContext{},
			WithShareMode(ShareExclusive),
			WithProtocol(ProtocolT1),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := ctx.readers[0], "Test"; got != want {
			t.Fatalf("ctx.readers[0] = %q, want %q", got, want)
		}
	})
}

func TestContextRelease(t *testing.T) {
	t.Run("Error from Release", func(t *testing.T) {
		ctx, err := newContext(&mockContext{
			release: func() error {
				return scard.ErrUnknownError
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := ctx.Release(); err != scard.ErrUnknownError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		ctx, err := newContext(&mockContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := ctx.Release(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestContextReaders(t *testing.T) {
	readers := []string{"r1", "r2"}

	ctx := &Context{readers: readers}

	if got, want := ctx.Readers(), readers; !stringsEqual(got, want) {
		t.Fatalf("ctx.Readers() = %v, want %v", got, want)
	}
}

func TestContextConnect(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctx, err := newContext(&mockContext{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = ctx.connect("Test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestContextWaitUntilCardPresent(t *testing.T) {
	t.Run("Error from GetStatusChange", func(t *testing.T) {
		ctx, err := newContext(&mockContext{
			getStatusChange: func(rs []scard.ReaderState, timeout time.Duration) error {
				return scard.ErrUnknownError
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := ctx.waitUntilCardPresent(); err != scard.ErrUnknownError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		ctx, err := newContext(&mockContext{
			getStatusChange: getStatusChangeFunc(scard.StatePresent),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reader, err := ctx.waitUntilCardPresent()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := reader, "Test"; got != want {
			t.Fatalf("reader = %q, want %q", got, want)
		}
	})
}

func TestContextWaitUntilCardRelease(t *testing.T) {
	t.Run("Error from GetStatusChange", func(t *testing.T) {
		ctx, err := newContext(&mockContext{
			getStatusChange: func(rs []scard.ReaderState, timeout time.Duration) error {
				return scard.ErrUnknownError
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := ctx.waitUntilCardRelease("Test"); err != scard.ErrUnknownError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		ctx, err := newContext(&mockContext{
			getStatusChange: getStatusChangeFunc(scard.StateEmpty),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := ctx.waitUntilCardRelease("Test"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

type mockContext struct {
	release         func() error
	isValid         func() (bool, error)
	listReaders     func() ([]string, error)
	connect         func(string, scard.ShareMode, scard.Protocol) (*scard.Card, error)
	getStatusChange func([]scard.ReaderState, time.Duration) error
}

func (ctx *mockContext) Release() error {
	if ctx.release != nil {
		return ctx.release()
	}

	return nil
}

func (ctx *mockContext) IsValid() (bool, error) {
	if ctx.isValid != nil {
		return ctx.isValid()
	}

	return true, nil
}

func (ctx *mockContext) ListReaders() ([]string, error) {
	if ctx.listReaders != nil {
		return ctx.listReaders()
	}

	return []string{"Test"}, nil
}

func (ctx *mockContext) Connect(reader string, shareMode scard.ShareMode, protocol scard.Protocol) (*scard.Card, error) {
	if ctx.connect != nil {
		return ctx.connect(reader, shareMode, protocol)
	}

	return &scard.Card{}, nil
}

func (ctx *mockContext) GetStatusChange(rs []scard.ReaderState, timeout time.Duration) error {
	if ctx.getStatusChange != nil {
		return ctx.getStatusChange(rs, timeout)
	}

	for i := range rs {
		rs[i].EventState = scard.StatePresent
	}

	return nil
}

func getStatusChangeFunc(sf scard.StateFlag) func([]scard.ReaderState, time.Duration) error {
	return func(rs []scard.ReaderState, timeout time.Duration) error {
		for i := range rs {
			rs[i].EventState = sf
		}

		return nil
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, s := range a {
		if s != b[i] {
			return false
		}
	}

	return true
}
