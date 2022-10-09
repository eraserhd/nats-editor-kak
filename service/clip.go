package service

import (
	"log"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type clipChangedAction struct {
	msg              *nats.Msg
	runKakouneScript func(cmd kakoune.Command) error
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(kakoune.Command{})
}
