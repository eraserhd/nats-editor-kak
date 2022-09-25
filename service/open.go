package service

import (
	"bytes"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"
)

var templ = template.Must(template.New("script").Parse(`
        edit -existing {{.QuotedFilename}}
`))

type Script struct {
	QuotedFilename string
}

func (s Script) String() string {
	buf := &bytes.Buffer{}
	_ = templ.Execute(buf, s)
	return buf.String()
}

type OpenCmd struct {
	Session   string
	Script    Script
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
			QuotedFilename: quote(u.Path),
		},
	}
	return result
}
