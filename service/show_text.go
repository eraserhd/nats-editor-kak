package service

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showText struct {
	Client string
	Text   string
	Base   string
}

var showTextTempl = template.Must(template.New("script").Parse(`
  kakoune-pluggo-show-text {{.Client}} {{.Text}} {{.Base}}
`))

func (s *showText) String() string {
	buf := &bytes.Buffer{}
	showTextTempl.Execute(buf, s)
	return buf.String()
}

func executeShowText(a *action) {
	a.log("info", fmt.Sprintf("received text to show: %q", string(a.msg.Data)))
	cmd := kakoune.Command{
		Session: a.kakouneSession,
		Script: &showText{
			Client: "%opt{jumpclient}",
			Text:   kakoune.Quote(string(a.msg.Data)),
			Base:   kakoune.Quote(a.msg.Header.Get("Base")),
		},
	}
	a.runKakouneScript(cmd)

	reply := nats.NewMsg(a.msg.Reply)
	reply.Data = []byte("ok")
	a.publish(reply)
}
