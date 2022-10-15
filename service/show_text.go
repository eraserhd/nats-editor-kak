package service

import (
	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showTextAction struct {
	kakouneSession   string
	msg              *nats.Msg
	publish          func(msg *nats.Msg) error
	runKakouneScript func(cmd kakoune.Command) error
}

func (a *showTextAction) Execute() {
}
