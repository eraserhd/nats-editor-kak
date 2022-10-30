package service

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showText struct {
	Client string
	Text   string
	Base   string
}

func (s *showText) String() string {
	return fmt.Sprintf("pluggo-show-text %s %s %s", s.Client, s.Text, s.Base)
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
