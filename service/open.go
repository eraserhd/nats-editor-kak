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

func (s *OpenScript) String() string {
	buf := &bytes.Buffer{}
	_ = templ.Execute(buf, s)
	return buf.String()
}

type OpenCommand struct {
	Session string
	Script  fmt.Stringer
}

type openAction struct {
	msg              *nats.Msg
	publish          func(msg *nats.Msg) error
	runKakouneScript func(cmd OpenCommand) error
}

func (a *openAction) makeOpenScript() OpenCommand {
	u, _ := url.Parse(string(a.msg.Data))
        openScript := &OpenScript{
		Client:         "%opt{jumpclient}",
		QuotedFilename: quote(u.Path),
		Selection: fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 1, Column: 1},
			End:   fragment.LineAndColumn{Line: 1, Column: 1},
		},
		FixupKeys: "''",
	}
	result := OpenCommand{
		Session: "kakoune",
		Script: openScript,
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
				openScript.FixupKeys = "'<a-L>'"
			}
			openScript.Selection = frag
		}
	}
	if s := a.msg.Header.Get("Session"); s != "" {
		result.Session = s
	}
	if w := a.msg.Header.Get("Window"); w != "" {
		openScript.Client = quote(w)
	}
	return result
}

func (a *openAction) Execute() {
	log.Printf("recieved %q", string(a.msg.Data))

	cmd := a.makeOpenScript()
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
