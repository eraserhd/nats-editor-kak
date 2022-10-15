package service

import (
	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type Service struct {
	kakouneSession string
}

type Action interface {
	Execute()
}

func New(kakouneSession string) (*Service, error) {
	return &Service{kakouneSession: kakouneSession}, nil
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

	textCh := make(chan *nats.Msg, 32)
	textSub, err := nc.ChanSubscribe("cmd.show.data.text", textCh)
	if err != nil {
		return err
	}
	defer textSub.Drain()

	clipCh := make(chan *nats.Msg, 32)
	clipSub, err := nc.ChanSubscribe("event.changed.clipboard", clipCh)
	if err != nil {
		return err
	}
	defer clipSub.Drain()

	for {
		var action Action
		select {
		case msg := <-fileCh:
			action = &showFileURLAction{
				kakouneSession: s.kakouneSession,
				msg:            msg,
				publish: func(msg *nats.Msg) error {
					return nc.PublishMsg(msg)
				},
				runKakouneScript: kakoune.Run,
			}
		case msg := <-textCh:
			action = &showTextAction{
				kakouneSession: s.kakouneSession,
				msg:            msg,
				publish: func(msg *nats.Msg) error {
					return nc.PublishMsg(msg)
				},
				runKakouneScript: kakoune.Run,
			}
		case msg := <-clipCh:
			action = &clipChangedAction{
				kakouneSession:   s.kakouneSession,
				msg:              msg,
				runKakouneScript: kakoune.Run,
			}
		}
		action.Execute()
	}
}
