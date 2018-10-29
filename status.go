package acr122u

import "github.com/ebfe/scard"

// Status contains the status of a card
type Status struct {
	Reader         string
	State          uint32
	ActiveProtocol uint32
	Atr            []byte
}

func newStatus(scs *scard.CardStatus) (Status, error) {
	if scs == nil {
		return Status{}, scard.ErrUnknownCard
	}

	return Status{
		Reader:         scs.Reader,
		State:          uint32(scs.State),
		ActiveProtocol: uint32(scs.ActiveProtocol),
		Atr:            scs.Atr,
	}, nil
}
