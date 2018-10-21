# acr122u

[![Build Status](https://travis-ci.org/peterhellberg/acr122u.svg?branch=master)](https://travis-ci.org/peterhellberg/acr122u)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/acr122u)](https://goreportcard.com/report/github.com/peterhellberg/acr122u)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/acr122u)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/acr122u#license-mit)

## Installation

    go get -u github.com/peterhellberg/acr122u

<img src="http://downloads.acs.com.hk/product-website-image/acr38-image.jpg" align="right" width="230" height="230">

## Dependencies

 - <https://www.acs.com.hk/en/products/3/acr122u-usb-nfc-reader/> - ACR122U USB NFC Reader
 - <https://pcsclite.apdu.fr/> - Middleware to access a smart card using SCard API (PC/SC)
 - <https://github.com/ebfe/scard> - Go bindings to the PC/SC API

 Under macOS `pcsc-lite` can be installed using homebrew: `brew install pcsc-lite`

## License (MIT)

Copyright (c) 2018 [Peter Hellberg](https://c7.se/)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
