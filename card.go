package acr122u

import "github.com/ebfe/scard"

var cmdReadTagID = []byte{0xFF, 0xCA, 0x00, 0x00, 0x04}

type scardCard interface {
	Transmit([]byte) ([]byte, error)
	Status() (*scard.CardStatus, error)
	Disconnect(scard.Disposition) error
}

// Card contains a ACR122U card
type Card struct {
	scard scardCard
}

func newCard(sc scardCard) *Card {
	return &Card{scard: sc}
}

// ReadTagID returns the tag ID for the card
func (c *Card) ReadTagID() ([]byte, error) {
	return c.scard.Transmit(cmdReadTagID)
}
