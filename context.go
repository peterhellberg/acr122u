package acr122u

import "github.com/ebfe/scard"

var scardEstablishContext = scard.EstablishContext

// Context for ACR122U readers
type Context struct {
	context   scardContext
	readers   []string
	shareMode ShareMode
	protocol  Protocol
}

// EstablishContext creates a ACR122U context
func EstablishContext(options ...Option) (*Context, error) {
	sctx, err := scardEstablishContext()
	if err != nil {
		return nil, err
	}

	return newContext(sctx, options...)
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
		context:   sctx,
		readers:   readers,
		shareMode: ShareShared,
		protocol:  ProtocolAny,
	}

	for _, option := range options {
		option(ctx)
	}

	return ctx, nil
}

// Release should be called when the context is not needed anymore
func (ctx *Context) Release() error {
	return ctx.context.Release()
}

// Readers returns a list of readers
func (ctx *Context) Readers() []string {
	return ctx.readers
}

// ServeFunc uses the provided HandlerFunc as a Handler
func (ctx *Context) ServeFunc(hf HandlerFunc) error {
	return ctx.Serve(hf)
}

// Serve cards being swiped using the provided Handler
func (ctx *Context) Serve(h Handler) error {
	for {
		ctx.serve(h)
	}
}

func (ctx *Context) serve(h Handler) error {
	reader, err := ctx.waitUntilCardPresent()
	if err != nil {
		return err
	}

	c, err := ctx.connect(reader)
	if err != nil {
		return err
	}

	if c.uid, err = c.getUID(); err == nil {
		h.ServeCard(c)
	} else {
		return err
	}

	return ctx.waitUntilCardRelease(reader)
}

func (ctx *Context) connect(reader string) (*card, error) {
	sc, err := ctx.context.Connect(reader,
		scard.ShareMode(ctx.shareMode),
		scard.Protocol(ctx.protocol),
	)
	if err != nil {
		return nil, err
	}

	return newCard(reader, sc), nil
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
