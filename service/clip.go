package service

import (
	"log"

	"github.com/nats-io/nats.go"
)

type clipChangedAction struct {
	msg              *nats.Msg
	runKakouneScript func(cmd KakouneCommand)
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(KakouneCommand{})
}
