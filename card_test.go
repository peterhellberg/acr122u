package acr122u

import (
	"bytes"
	"testing"

	"github.com/ebfe/scard"
)

func TestCardReadTagID(t *testing.T) {
	tagID := []byte{0x83, 0xfb, 0x58, 0x24, 0x90}

	c := newCard(&mockCard{
		transmit: func(cmd []byte) ([]byte, error) {
			if !bytes.Equal(cmd, cmdReadTagID) {
				t.Fatal("cmd != commandReadTagID")
			}

			return tagID, nil
		},
	})

	got, err := c.ReadTagID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, tagID) {
		t.Fatalf("%#v != %#v", got, tagID)
	}
}

type mockCard struct {
	transmit   func([]byte) ([]byte, error)
	status     func() (*scard.CardStatus, error)
	disconnect func(scard.Disposition) error
}

func (c *mockCard) Transmit(cmd []byte) ([]byte, error) {
	return c.transmit(cmd)
}

func (c *mockCard) Status() (*scard.CardStatus, error) {
	return c.status()
}

func (c *mockCard) Disconnect(d scard.Disposition) error {
	return c.disconnect(d)
}
