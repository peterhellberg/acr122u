package main

import (
	"log"

	"github.com/peterhellberg/acr122u"
)

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		log.Println(err)
		return
	}
	defer ctx.Release()

	for {
		if err := ctx.WhenCardPresent(handleCard); err != nil {
			log.Println(err)
			break
		}
	}
}

func handleCard(c *acr122u.Card) error {
	tagID, err := c.ReadTagID()
	if err != nil {
		return err
	}

	log.Printf("%x\n", tagID)

	return nil
}
