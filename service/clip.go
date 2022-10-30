package service

import (
	"fmt"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type SetDquoteRegister struct {
	Value string
}

func (s *SetDquoteRegister) String() string {
	return fmt.Sprintf("pluggo-clip %s", s.Value)
}

func executeClipChanged(a *action) {
	a.log("info", fmt.Sprintf("clipboard changed: %q", string(a.msg.Data)))
	a.runKakouneScript(kakoune.Command{
		Session: a.kakouneSession,
		Script: &SetDquoteRegister{
			Value: kakoune.Quote(string(a.msg.Data)),
		},
	})
}
