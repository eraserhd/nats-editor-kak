package service

import (
	"bytes"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"
)

type Script struct {
	Client         string
	QuotedFilename string
}

var templ = template.Must(template.New("script").Parse(`
  evalutate-commands -try-client {{.Client}} %{
    try %{
      edit -existing {{.QuotedFilename}}
      try focus
    } catch %{
      echo -markup "{Error}%val{error}"
      echo -debug "%val{error}"
    }
  }
`))

func (s Script) String() string {
	buf := &bytes.Buffer{}
	_ = templ.Execute(buf, s)
	return buf.String()
}

type OpenCmd struct {
	Session string
	Script  Script
}

func quote(s string) string {
	result := "'"
	for _, ch := range s {
		if ch == '\'' {
			result += "'"
		}
		result += string(ch)
	}
	return result + "'"
}

func (s *Service) OpenCommand(msg *nats.Msg) OpenCmd {
	u, _ := url.Parse(string(msg.Data))
	result := OpenCmd{
		Session: msg.Header.Get("Session"),
		Script: Script{
			Client:         "%opt{jumpclient}",
			QuotedFilename: quote(u.Path),
		},
	}
	if w := msg.Header.Get("Window"); w != "" {
        	result.Script.Client = quote(w)
	}
	return result
}
