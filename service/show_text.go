package service

import (
	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type action struct {
	kakouneSession   string
	msg              *nats.Msg
	publish          func(msg *nats.Msg) error
	runKakouneScript func(cmd kakoune.Command) error
}

type showTextAction action

func (a *showTextAction) Execute() {
}
