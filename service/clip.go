package service

import (
	"log"

	"github.com/nats-io/nats.go"
)

type ChangeDquoteRegister struct {
}

type clipChangedAction struct {
	msg              *nats.Msg
	runKakouneScript func(cmd ChangeDquoteRegister)
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(ChangeDquoteRegister{})
}
