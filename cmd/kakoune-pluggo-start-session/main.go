package main

import (
	"flag"
	"os"
	"path"
	"text/template"
)

type Params struct {
        BinPath string
	Session string
}

var initScript = template.Must(template.New("init").Parse(`
nop %sh{
    {{.BinPath}}/kakoune-pluggo-daemon {{.Session}} </dev/null >/dev/null 2>&1 &
}
`))

func binPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return path.Dir(exe)
}

func reportError(s string) {
	panic(s)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		reportError("Syntax: kakoune-pluggo-start-session SESSION-NAME")
	}
	params := Params{
		BinPath: binPath(),
		Session: args[0],
	}
	if err := initScript.Execute(os.Stdout, params); err != nil {
		reportError(err.Error())
	}
}
