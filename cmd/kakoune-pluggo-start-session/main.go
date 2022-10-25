package main

import (
	_ "embed"
	"os"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
	"github.com/plugbench/kakoune-pluggo/service"
)

type ScriptParams struct {
	BinPath string
}

//go:embed start-session.kak
var templateSource string

var scriptTemplate = template.Must(template.New("start-session.kak").Parse(templateSource))

func main() {
	params := ScriptParams{
		BinPath: service.BinPath(),
	}
	if err := scriptTemplate.Execute(os.Stdout, params); err != nil {
		kakoune.Fail(err.Error())
	}
}
