package main

import (
	"errors"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

func sendCommand(subject, data string) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	req := nats.NewMsg(subject)
	req.Data = []byte(data)

	response, err := nc.RequestMsg(req, 3*time.Second)
	if err != nil {
		return err
	}

	if string(response.Data) != "ok" {
		return errors.New(string(response.Data))
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		kakoune.Fail("Syntax is kakoune-pluggo-command SUBJECT DATA")
	}
	if err := sendCommand(os.Args[1], os.Args[2]); err != nil {
		kakoune.Fail(err.Error())
	}
}
