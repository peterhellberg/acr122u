package main

import "github.com/peterhellberg/acr122u"

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	h := &handler{acr122u.StdoutLogger()}

	ctx.Serve(h)
}

type handler struct {
	acr122u.Logger
}

func (h *handler) ServeCard(c acr122u.Card) {
	h.Printf("%x\n", c.UID())
}
