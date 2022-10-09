package service

import (
	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type Service struct {
}

func New() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	fileCh := make(chan *nats.Msg, 32)
	fileSub, err := nc.ChanSubscribe("cmd.show.url.file", fileCh)
	if err != nil {
		return err
	}
	defer fileSub.Drain()

	clipCh := make(chan *nats.Msg, 32)
	clipSub, err := nc.ChanSubscribe("event.changed.clipboard", clipCh)
	if err != nil {
		return err
	}
	defer clipSub.Drain()

	for {
		select {
		case msg := <-fileCh:
			action := openAction{
				msg: msg,
				publish: func(msg *nats.Msg) error {
					return nc.PublishMsg(msg)
				},
				runKakouneScript: kakoune.Run,
			}
			action.Execute()
		case msg := <-clipCh:
			action := clipChangedAction{
				msg: msg,
				runKakouneScript: kakoune.Run,
			}
			action.Execute()
		}
	}
}
