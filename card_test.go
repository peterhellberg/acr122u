package acr122u

import (
	"bytes"
	"testing"

	"github.com/ebfe/scard"
)

func TestNewCard(t *testing.T) {
	m := &mockCard{}
	c := newCard("", m)

	if got, want := c.scard.(*mockCard), m; got != want {
		t.Fatalf("c.scard = %v, want %v", got, want)
	}
}

func TestCardReader(t *testing.T) {
	r := "test-reader"
	c := newCard(r, nil)

	if got, want := c.Reader(), r; got != want {
		t.Fatalf("c.Reader() = %q, want %q", got, want)
	}
}

func TestCardStatus(t *testing.T) {
	t.Run("Error from Status", func(t *testing.T) {
		c := statusCard(func() (*scard.CardStatus, error) {
			return nil, scard.ErrUnknownError
		})

		if _, err := c.Status(); err != scard.ErrUnknownError {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		c := statusCard(func() (*scard.CardStatus, error) {
			return &scard.CardStatus{Reader: "Test"}, nil
		})

		s, err := c.Status()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := s.Reader, "Test"; got != want {
			t.Fatalf("s.Reader = %q, want %q", got, want)
		}
	})
}

func TestCardUID(t *testing.T) {
	c := &card{uid: testUID}

	if got := c.UID(); !bytes.Equal(got, testUID) {
		t.Fatalf("c.UID() = %#v, want %#v", got, testUID)
	}
}

func TestCardGetUID(t *testing.T) {
	c := transmitCard(func(cmd []byte) ([]byte, error) {
		if !bytes.Equal(cmd, cmdGetUID) {
			t.Fatalf("cmd = %v, want %v", cmd, cmdGetUID)
		}

		return testUID, nil
	})

	got, err := c.getUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, testUID) {
		t.Fatalf("%#v != %#v", got, testUID)
	}
}

var testUID = []byte{0x83, 0xfb, 0x58, 0x24, 0x90}

type mockCard struct {
	transmit func([]byte) ([]byte, error)
	status   func() (*scard.CardStatus, error)
}

func (c *mockCard) Transmit(cmd []byte) ([]byte, error) {
	return c.transmit(cmd)
}

func (c *mockCard) Status() (*scard.CardStatus, error) {
	return c.status()
}

func transmitCard(t func(cmd []byte) ([]byte, error)) *card {
	return newCard("", &mockCard{transmit: t})
}

func statusCard(s func() (*scard.CardStatus, error)) *card {
	return newCard("", &mockCard{status: s})
}
