package service

import (
	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showText struct {
	Client string
}

func (s *showText) String() string {
	return "" //FIXME:
}

func executeShowText(act *action) {
	cmd := kakoune.Command{
		Session: act.kakouneSession,
		Script: &showText{
			Client: "%opt{jumpclient}",
		},
	}
	act.runKakouneScript(cmd)
}
