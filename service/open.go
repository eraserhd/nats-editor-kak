package service

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/service/fragment"
)

type openAction struct {
	msg              *nats.Msg
	publish          func(msg *nats.Msg) error
	runKakouneScript func(cmd OpenCommand) error
}

func (a *openAction) Execute() {
	log.Printf("recieved %q", string(a.msg.Data))

	cmd := openCommand(a.msg)
	if err := a.runKakouneScript(cmd); err != nil {
		log.Print(err)
		reply := nats.NewMsg(a.msg.Reply)
		reply.Data = []byte(fmt.Sprintf("ERROR: %s", err.Error()))
		if err := a.publish(reply); err != nil {
			log.Printf("error responding: %v", err)
		}
		return
	}

	reply := nats.NewMsg(a.msg.Reply)
	reply.Data = []byte("ok")
	if err := a.publish(reply); err != nil {
		log.Printf("error replying ok: %v", err)
		return
	}
}

type OpenScript struct {
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

func (s OpenScript) String() string {
	buf := &bytes.Buffer{}
	_ = templ.Execute(buf, s)
	return buf.String()
}

type OpenCommand struct {
	Session string
	Script  OpenScript
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

func openCommand(msg *nats.Msg) OpenCommand {
	u, _ := url.Parse(string(msg.Data))
	result := OpenCommand{
		Session: "kakoune",
		Script: OpenScript{
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
			result.Script.Selection = frag
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
