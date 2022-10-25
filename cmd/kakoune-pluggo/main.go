package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
	"github.com/plugbench/kakoune-pluggo/service"
)

const help = `
Syntax:
  kakoune-pluggo SUBCOMMAND [OPTIONS...]

Subcommands:
  start-session       Print Kakoune initialization script and exit.  Intended to be invoked as
                      "evaluate-commands %sh{kakoune-pluggo start-session}".
`

type ScriptParams struct {
	BinPath string
}

//go:embed start-session.kak
var templateSource string

var scriptTemplate = template.Must(template.New("start-session.kak").Parse(templateSource))

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start-session":
		params := ScriptParams{
			BinPath: service.BinPath(),
		}
		if err := scriptTemplate.Execute(os.Stdout, params); err != nil {
			kakoune.Fail(err.Error())
		}

	default:
		fmt.Println(help)
		os.Exit(1)
	}
}
