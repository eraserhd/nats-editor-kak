package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"


	"github.com/plugbench/kakoune-pluggo/kakoune"
	knats "github.com/plugbench/kakoune-pluggo/nats"
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
	Nats      knats.Config
}

//go:embed start-session.kak
var templateSource string

var scriptTemplate = template.Must(template.New("start-session.kak").Parse(templateSource))

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(help)
		os.Exit(1)
	}
	nats := knats.LoadConfigFromEnvironment()
	switch os.Args[1] {
	case "command":
		if len(os.Args) != 4 {
			kakoune.Fail("wrong argument count, see help")
		}
		if err := nats.SendCommand(os.Args[2], os.Args[3]); err != nil {
			kakoune.Fail(err.Error())
		}

	case "daemon":
		if len(os.Args) != 3 {
			kakoune.Fail("wrong argument count, see help")
		}
		session := os.Args[2]
		es, err := service.New(nats, session)
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
		if err := nats.SendEvent(os.Args[2], os.Args[3]); err != nil {
			kakoune.Fail(err.Error())
		}

	case "start-session":
		params := ScriptParams{
			PluggoBin: service.PluggoBin(),
			Nats:      nats,
		}
		if err := scriptTemplate.Execute(os.Stdout, params); err != nil {
			kakoune.Fail(err.Error())
		}

	default:
		fmt.Printf(help)
		os.Exit(1)
	}
}
