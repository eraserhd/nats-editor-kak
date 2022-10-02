package service

import (
	"bytes"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/service/fragment"
)

type Script struct {
	Client         string
	QuotedFilename string
	Selection      fragment.LineAndColumnSelection
	FixupKeys      string
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
			Selection: fragment.LineAndColumnSelection{
				Start: fragment.LineAndColumn{Line: 1, Column: 1},
				End:   fragment.LineAndColumn{Line: 1, Column: 1},
			},
			FixupKeys: "''",
		},
	}
	if frag, err := fragment.ParseRFC5147FragmentIdentifier(u.Fragment); err == nil {
		if frag, ok := frag.(fragment.LineAndColumnSelection); ok {
			if frag.Start == frag.End {
				// Kakoune doesn't do zero-width selections.
				frag.End.Column++
			}
        		// Adjust for Kakoune's 1-based indexing
			frag.Start.Line++
			frag.Start.Column++
			frag.End.Line++
			if frag.End.Column == 0 {
				// Kakoune can't select up to the zero-width point at BOL, so if we are trying
				// to do so, select up to the previous line and extend to EOL with <a-L>
				frag.End.Line--
				frag.End.Column = 1
				result.Script.FixupKeys = "'<a-L>'"
			}
			result.Script.Selection = fragment.LineAndColumnSelection{
				Start: frag.Start,
				End:   frag.End,
			}
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
