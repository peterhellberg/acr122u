/*

Package acr122u is a library for the ACR122U USB NFC Reader

Requirements

ACR122U USB NFC Reader https://www.acs.com.hk/en/products/3/acr122u-usb-nfc-reader

Middleware to access a smart card using SCard API (PC/SC)  https://pcsclite.apdu.fr

    Under macOS pcsc-lite can be installed using homebrew: brew install pcsc-lite

The Go bindings to the PC/SC API https://github.com/ebfe/scard

Installation

You can install the acr122u package using go get

    go get -u github.com/peterhellberg/acr122u

Usage

A minimal usage example

    package main

    import (
    	"fmt"

    	"github.com/peterhellberg/acr122u"
    )

    func main() {
    	ctx, err := acr122u.EstablishContext()
    	if err != nil {
    		panic(err)
    	}

    	ctx.ServeFunc(func(c acr122u.Card) {
    		fmt.Printf("%x\n", c.UID())
    	})
    }

*/
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

// Response codes
var (
	rcOperationSuccess = []byte{0x90, 0x00}
	rcOperationFailed  = []byte{0x63, 0x00}
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
}
