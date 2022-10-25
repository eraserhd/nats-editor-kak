package service

import (
	"fmt"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showText struct {
	Client string
}

func (s *showText) String() string {
	return "" //FIXME:
}

func executeShowText(act *action) {
	act.log("info", fmt.Sprintf("received text to show: %q", string(act.msg.Data)))
	cmd := kakoune.Command{
		Session: act.kakouneSession,
		Script: &showText{
			Client: "%opt{jumpclient}",
		},
	}
	act.runKakouneScript(cmd)
}
