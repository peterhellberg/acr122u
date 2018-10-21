package main

import (
	"log"
	"os"

	"github.com/peterhellberg/acr122u"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	cmd, err := newCmd(logger)
	if err != nil {
		logger.Printf("%v\n", err)
		return
	}
	defer cmd.Release()

	cmd.Run()
}

type Cmd struct {
	*acr122u.Context
	*log.Logger
}

func newCmd(logger *log.Logger) (*Cmd, error) {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		return nil, err
	}

	return &Cmd{Context: ctx, Logger: logger}, nil
}

func (cmd *Cmd) HandleCard(c *acr122u.Card) error {
	tagID, err := c.ReadTagID()
	if err != nil {
		return err
	}

	cmd.log("%x", tagID)

	return nil
}

func (cmd *Cmd) Run() {
	for {
		if err := cmd.WhenCardPresent(cmd.HandleCard); err != nil {
			cmd.Println(err)
			break
		}
	}
}

func (cmd *Cmd) log(format string, v ...interface{}) {
	cmd.Printf(format+"\n", v...)
}
