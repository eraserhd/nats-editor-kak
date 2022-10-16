package service

import (
	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type Service struct {
	kakouneSession string
}

type action struct {
	kakouneSession   string
	msg              *nats.Msg
	publish          func(msg *nats.Msg) error
	runKakouneScript func(cmd kakoune.Command) error
	execute          func(a *action)
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
		act := action{
			kakouneSession: s.kakouneSession,
			publish: func(msg *nats.Msg) error {
				return nc.PublishMsg(msg)
			},
			runKakouneScript: kakoune.Run,
		}
		select {
		case act.msg = <-fileCh:
			act.execute = executeShowFileURL
		case act.msg = <-textCh:
			act.execute = executeShowText
		case act.msg = <-clipCh:
			act.execute = executeClipChanged
		}
		act.execute(&act)
	}
}
