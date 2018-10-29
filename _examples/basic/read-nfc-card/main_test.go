package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/peterhellberg/acr122u"
)

func TestHandlerServeCard(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	h := &handler{log.New(buf, "", 0)}

	h.ServeCard(&card{})

	if got, want := buf.String(), "010203\n"; got != want {
		t.Fatalf("buf.String() = %q, want %q", got, want)
	}
}

type card struct{}

func (c *card) Reader() string {
	return ""
}

func (c *card) Status() (acr122u.Status, error) {
	return acr122u.Status{}, nil
}

func (c *card) UID() []byte {
	return []byte{0x1, 0x2, 0x3}
}
