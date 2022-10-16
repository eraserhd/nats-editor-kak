package service

import (
	"github.com/plugbench/kakoune-pluggo/kakoune"
)

func executeShowText(act *action) {
	cmd := kakoune.Command{
		Session: act.kakouneSession,
	}
	act.runKakouneScript(cmd)
}
