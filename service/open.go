package service

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
	"github.com/plugbench/kakoune-pluggo/service/fragment"
)

type OpenFile struct {
	Client         string
	QuotedFilename string
	Selection      fragment.LineAndColumnSelection
	FixupKeys      string
}

var openFileTempl = template.Must(template.New("script").Parse(`
 pluggo-open {{.Client}} {{.QuotedFilename}} {{.Selection.Start.Line}}.{{.Selection.Start.Column}},{{.Selection.End.Line}}.{{.Selection.End.Column}} {{.FixupKeys}}
`))

func (s *OpenFile) String() string {
	buf := &bytes.Buffer{}
	_ = openFileTempl.Execute(buf, s)
	return buf.String()
}

func makeOpenScript(a *action) kakoune.Command {
	u, _ := url.Parse(string(a.msg.Data))
	openScript := &OpenFile{
		Client:         "%opt{jumpclient}",
		QuotedFilename: kakoune.Quote(u.Path),
		Selection: fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 1, Column: 1},
			End:   fragment.LineAndColumn{Line: 1, Column: 1},
		},
		FixupKeys: "''",
	}
	result := kakoune.Command{
		Session: a.kakouneSession,
		Script:  openScript,
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
	if w := a.msg.Header.Get("Window"); w != "" {
		openScript.Client = kakoune.Quote(w)
	}
	return result
}

func executeShowFileURL(a *action) {
	a.log("info", fmt.Sprintf("recieved %q", string(a.msg.Data)))

	cmd := makeOpenScript(a)
	if err := a.runKakouneScript(cmd); err != nil {
		a.log("error", fmt.Sprintf("error running open script: %v", err))
		reply := nats.NewMsg(a.msg.Reply)
		reply.Data = []byte(fmt.Sprintf("ERROR: %s", err.Error()))
		if err := a.publish(reply); err != nil {
			a.log("error", fmt.Sprintf("error responding: %v", err))
		}
		return
	}

	reply := nats.NewMsg(a.msg.Reply)
	reply.Data = []byte("ok")
	if err := a.publish(reply); err != nil {
		a.log("error", fmt.Sprintf("error replying ok: %v", err))
		return
	}
}
