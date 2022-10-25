package main

import (
	"os"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

func sendEvent(subject, data string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	req := nats.NewMsg(subject)
	req.Data = []byte(data)

	if err := nc.PublishMsg(req); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		kakoune.Fail("Syntax is kakoune-pluggo-event SUBJECT DATA")
	}
	if err := sendEvent(os.Args[1], os.Args[2]); err != nil {
		kakoune.Fail(err.Error())
	}
}

