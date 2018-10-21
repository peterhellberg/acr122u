package acr122u

import (
	"time"

	"github.com/ebfe/scard"
)

var (
	scardEstablishContext = scard.EstablishContext
	cmdBuzzerEnable       = []byte{0xFF, 0x00, 0x52, 0xFF, 0x00}
	cmdBuzzerDisable      = []byte{0xFF, 0x00, 0x52, 0x00, 0x00}
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

type scardContext interface {
	Release() error
	IsValid() (bool, error)
	ListReaders() ([]string, error)
	Connect(string, scard.ShareMode, scard.Protocol) (*scard.Card, error)
	GetStatusChange(readerStates []scard.ReaderState, timeout time.Duration) error
}

// EstablishContext creates a ACR122U context
func EstablishContext(options ...Option) (*Context, error) {
	sctx, err := scardEstablishContext()
	if err != nil {
		return nil, err
	}

	return newContext(sctx, options...)
}

func newContext(sctx scardContext, options ...Option) (*Context, error) {
	if _, err := sctx.IsValid(); err != nil {
		return nil, err
	}

	readers, err := sctx.ListReaders()
	if err != nil {
		return nil, err
	}

	if len(readers) == 0 {
		return nil, scard.ErrNoReadersAvailable
	}

	ctx := &Context{
		context:        sctx,
		readers:        readers,
		shareMode:      ShareShared,
		protocol:       ProtocolAny,
		disabledBuzzer: true,
	}

	for _, option := range options {
		option(ctx)
	}

	return ctx, nil
}

// Option is the function type used to configure the context
type Option func(*Context)

// WithShareMode accepts Exclusive (0x1) or Shared mode (0x2)
func WithShareMode(sm ShareMode) Option {
	return func(ctx *Context) {
		ctx.shareMode = sm
	}
}

// WithProtocol accepts Undefined (0x0), T0 (0x1), T1 (0x2) or Any (T0|T1)
func WithProtocol(p Protocol) Option {
	return func(ctx *Context) {
		ctx.protocol = p
	}
}

// WithDisabledBuzzer is used to specify if the buzzer should be disabled or not
func WithDisabledBuzzer(b bool) Option {
	return func(ctx *Context) {
		ctx.disabledBuzzer = b
	}
}

// Context for ACR122U readers
type Context struct {
	context        scardContext
	readers        []string
	shareMode      ShareMode
	protocol       Protocol
	disabledBuzzer bool
}

// Release should be called when the context is not needed anymore
func (ctx *Context) Release() error {
	return ctx.context.Release()
}

func (ctx *Context) connect(reader string) (*Card, error) {
	sc, err := ctx.context.Connect(reader, scard.ShareMode(ctx.shareMode), scard.Protocol(ctx.protocol))
	if err != nil {
		return nil, err
	}

	return newCard(sc), nil
}

// WhenCardPresent accepts a function that is called each time a card is present
func (ctx *Context) WhenCardPresent(f func(*Card) error) error {
	reader, err := ctx.waitUntilCardPresent()
	if err != nil {
		return err
	}

	c, err := ctx.connect(reader)
	if err != nil {
		return err
	}

	// Disable the beep after the first card swipe
	if ctx.disabledBuzzer {
		c.scard.Transmit(cmdBuzzerDisable)
	} else {
		c.scard.Transmit(cmdBuzzerEnable)
	}

	if err := f(c); err != nil {
		return err
	}

	if err := c.scard.Disconnect(scard.ResetCard); err != nil {
		return err
	}

	if err := ctx.waitUntilCardRelease(reader); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) waitUntilCardPresent() (string, error) {
	rs := make([]scard.ReaderState, len(ctx.readers))

	for i := range rs {
		rs[i].Reader = ctx.readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return ctx.readers[i], nil
			}
			rs[i].CurrentState = rs[i].EventState
		}

		if err := ctx.context.GetStatusChange(rs, -1); err != nil {
			return "", err
		}
	}
}

func (ctx *Context) waitUntilCardRelease(reader string) error {
	rs := []scard.ReaderState{{
		Reader:       reader,
		CurrentState: scard.StatePresent,
	}}

	for {
		if rs[0].EventState&scard.StateEmpty != 0 {
			return nil
		}

		rs[0].CurrentState = rs[0].EventState

		if err := ctx.context.GetStatusChange(rs, -1); err != nil {
			return err
		}
	}
}
