package service

import (
	"log"

	"github.com/nats-io/nats.go"
)

type clipChangedAction struct {
	msg     *nats.Msg
	publish func(msg *nats.Msg) error
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
}
