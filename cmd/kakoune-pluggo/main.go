package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
	"github.com/plugbench/kakoune-pluggo/service"
)

const help = `
Syntax:
  kakoune-pluggo SUBCOMMAND [OPTIONS...]

Subcommands:
  command SUBJECT DATA  Send NATS command to SUBJECT, wait for and print a reply.

  daemon SESSION        Run the daemon for Kakoune SESSION.

  event SUBJECT DATA    Send NATS event to SUBJECT.

  start-session         Print Kakoune initialization script and exit.  Intended to be invoked as
                        "evaluate-commands %%sh{kakoune-pluggo start-session}".
`

type ScriptParams struct {
	PluggoBin string
}

//go:embed start-session.kak
var templateSource string

var scriptTemplate = template.Must(template.New("start-session.kak").Parse(templateSource))

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
	if len(os.Args) < 2 {
		fmt.Printf(help)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "command":
		if len(os.Args) != 4 {
			kakoune.Fail("wrong argument count, see help")
		}
		if err := sendCommand(os.Args[2], os.Args[3]); err != nil {
			kakoune.Fail(err.Error())
		}

	case "daemon":
		if len(os.Args) != 3 {
			kakoune.Fail("wrong argument count, see help")
		}
		session := os.Args[2]
		es, err := service.New(session)
		if err != nil {
			kakoune.Fail(err.Error())
		}
		if err := es.Run(); err != nil {
			kakoune.Fail(err.Error())
		}

	case "event":
		if len(os.Args) != 4 {
			kakoune.Fail("wrong argument count, see help")
		}
		if err := sendEvent(os.Args[2], os.Args[3]); err != nil {
			kakoune.Fail(err.Error())
		}

	case "start-session":
		params := ScriptParams{
			PluggoBin: service.PluggoBin(),
		}
		if err := scriptTemplate.Execute(os.Stdout, params); err != nil {
			kakoune.Fail(err.Error())
		}

	default:
		fmt.Printf(help)
		os.Exit(1)
	}
}
