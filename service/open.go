package service

import (
	"bytes"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/eraserhd/nats-editor-kak/service/fragment"
)

type Script struct {
	Client         string
	QuotedFilename string
	Selection      Selection
	FixupKeys      string
}

type Selection struct {
	Start, End Position
}

type Position struct {
	Line, Column int
}

var templ = template.Must(template.New("script").Parse(`
  evaluate-commands -try-client {{.Client}} %{
    try %{
      edit -existing {{.QuotedFilename}}
      select -codepoint {{.Selection.Start.Line}}.{{.Selection.Start.Column}},{{.Selection.End.Line}}.{{.Selection.End.Column}}
      {{ if ne .FixupKeys "" -}}
      execute-keys {{.FixupKeys}}
      {{- end }}
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
		Session: "kakoune",
		Script: Script{
			Client:         "%opt{jumpclient}",
			QuotedFilename: quote(u.Path),
			Selection: Selection{
				Start: Position{1, 1},
				End:   Position{1, 1},
			},
		},
	}
	if line, endLine, err := fragment.Parse(u.Fragment); err == nil {
		result.Script.Selection = Selection{
			Start: Position{int(line) + 1, 1},
			End:   Position{int(line) + 1, 1},
		}
		if line != endLine {
			result.Script.Selection.End.Line = int(endLine) - 1
        		result.Script.FixupKeys = "'x'" // FIXME: Only for line ranges
		}
	}
	if s := msg.Header.Get("Session"); s != "" {
		result.Session = s
	}
	if w := msg.Header.Get("Window"); w != "" {
		result.Script.Client = quote(w)
	}
	return result
}
