package main

import (
	"log"

	nats "github.com/nats-io/go-nats"
	acr122u "github.com/peterhellberg/acr122u"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := acr122u.EstablishContext()
	if err != nil {
		log.Fatal(err)
	}

	ctx.ServeFunc(func(c acr122u.Card) {
		nc.Publish("acr122u", c.UID())
	})
}
