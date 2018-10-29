package acr122u

import "github.com/ebfe/scard"

// Card represents a ACR122U card
type Card interface {
	// Reader returns the name of the reader used
	Reader() string

	// Status returns the card status
	Status() (Status, error)

	// UID returns the UID for the card
	UID() []byte
}

type card struct {
	uid    []byte
	reader string
	scard  scardCard
}

func newCard(reader string, sc scardCard) *card {
	return &card{reader: reader, scard: sc}
}

func (c *card) Reader() string {
	return c.reader
}

func (c *card) Status() (Status, error) {
	scs, err := c.scard.Status()
	if err != nil {
		return Status{}, err
	}

	return newStatus(scs)
}

func (c *card) UID() []byte {
	return c.uid
}

// transmit raw command to underlying scardCard
func (c *card) transmit(cmd []byte) ([]byte, error) {
	return c.scard.Transmit(cmd)
}

// getUID returns the UID for the card
func (c *card) getUID() ([]byte, error) {
	return c.transmit(cmdGetUID)
}

// disconnect the card
func (c *card) disconnect(d scard.Disposition) error {
	return c.scard.Disconnect(d)
}
