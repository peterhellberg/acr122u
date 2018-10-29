package main

import (
	"fmt"
	"runtime"

	nats "github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	nc.Subscribe("acr122u", func(m *nats.Msg) {
		fmt.Printf("Received UID: %x\n", m.Data)
	})

	runtime.Goexit()
}
