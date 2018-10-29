package main

import (
	"log"
	"os"

	"github.com/peterhellberg/acr122u"
)

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	h := &handler{log.New(os.Stdout, "", 0)}

	ctx.Serve(h)
}

type handler struct {
	acr122u.Logger
}

func (h *handler) ServeCard(c acr122u.Card) {
	h.Printf("%x\n", c.UID())
}
