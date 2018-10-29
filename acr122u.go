package acr122u

import (
	"time"

	"github.com/ebfe/scard"
)

// ShareMode is the share mode type
type ShareMode uint32

// Share modes
var (
	ShareExclusive ShareMode = 0x1
	ShareShared    ShareMode = 0x2
)

// Protocol is the protocol type
type Protocol uint32

// Protocols
var (
	ProtocolUndefined Protocol = 0x0
	ProtocolT0        Protocol = 0x1
	ProtocolT1        Protocol = 0x2
	ProtocolAny                = ProtocolT0 | ProtocolT1
)

// Commands that can be transmitted to a *scard.Card
var (
	cmdGetUID = []byte{0xFF, 0xCA, 0x00, 0x00, 0x04}
)

// scardContext is the interface used to communicate
// with one or more ACR122U USB NFC Readers.
type scardContext interface {
	Connect(string, scard.ShareMode, scard.Protocol) (*scard.Card, error)
	ListReaders() ([]string, error)
	Release() error
	IsValid() (bool, error)
	GetStatusChange(readerStates []scard.ReaderState, timeout time.Duration) error
}

// scardCard is the interface used by a *card to
// communicate with the underlying *scard.Card
type scardCard interface {
	Transmit([]byte) ([]byte, error)
	Status() (*scard.CardStatus, error)
	Disconnect(scard.Disposition) error
}
