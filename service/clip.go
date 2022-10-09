package service

import (
	"log"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type SetDquoteRegister struct {
	Value string
}

func (s *SetDquoteRegister) String() string {
	return ""
}

type clipChangedAction struct {
	msg              *nats.Msg
	runKakouneScript func(cmd kakoune.Command) error
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(kakoune.Command{
		Script: &SetDquoteRegister{
			Value: kakoune.Quote(string(a.msg.Data)),
		},
	})
}
