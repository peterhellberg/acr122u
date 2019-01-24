# acr122u

<img src="http://downloads.acs.com.hk/product-website-image/acr38-image.jpg" align="right" width="230" height="230">

[![Build Status](https://travis-ci.org/peterhellberg/acr122u.svg?branch=master)](https://travis-ci.org/peterhellberg/acr122u)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/acr122u)](https://goreportcard.com/report/github.com/peterhellberg/acr122u)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/acr122u)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/acr122u#license-mit)

This is a Go package for the ACR122U USB NFC Reader

## Requirements

 - <https://www.acs.com.hk/en/products/3/acr122u-usb-nfc-reader/> - ACR122U USB NFC Reader
 - <https://pcsclite.apdu.fr/> - Middleware to access a smart card using SCard API (PC/SC)
 - <https://github.com/ebfe/scard> - Go bindings to the PC/SC API

 Under macOS `pcsc-lite` can be installed using homebrew: `brew install pcsc-lite`

## Installation

    go get -u github.com/peterhellberg/acr122u

## Usage

### Minimal example

```go
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
```

### Using a struct that implements the `acr122u.Handler` interface

```go
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
```

### NATS

[NATS](https://nats.io/) is my favorite messaging system,
so let’s see how we can use it with this package:

#### Publish each UID to a NATS subject

```go
package main

import (
	nats "github.com/nats-io/go-nats"
	acr122u "github.com/peterhellberg/acr122u"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	ctx.ServeFunc(func(c acr122u.Card) {
		nc.Publish("acr122u", c.UID())
	})
}
```

#### Subscribe to the NATS subject

```go
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
```

## License (MIT)

Copyright (c) 2018-2019 [Peter Hellberg](https://c7.se)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
