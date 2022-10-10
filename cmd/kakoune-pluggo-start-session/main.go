package main

import (
	_ "embed"
	"os"
	"path"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type ScriptParams struct {
	BinPath string
	Session string
}

//go:embed start-session.kak
var templateSource string

var scriptTemplate = template.Must(template.New("start-session.kak").Parse(templateSource))

func binPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return path.Dir(exe)
}

func main() {
	if len(os.Args) != 2 {
		kakoune.Fail("Syntax: kakoune-pluggo-start-session SESSION-NAME")
	}
	params := ScriptParams{
		BinPath: binPath(),
		Session: os.Args[1],
	}
	if err := scriptTemplate.Execute(os.Stdout, params); err != nil {
		kakoune.Fail(err.Error())
	}
}
