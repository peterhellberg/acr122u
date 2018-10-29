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
