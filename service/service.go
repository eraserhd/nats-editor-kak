package service

import (
	"log"

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
}

func (a *action) dispatch() {
	switch a.msg.Subject {
	case "cmd.show.file.url":
		executeShowFileURL(a)
	case "cmd.show.data.text":
		executeShowText(a)
	case "event.changed.clipboard":
		executeClipChanged(a)
	default:
		log.Fatalf("do not know how to handle %s", a.msg.Subject)
	}
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
		case act.msg = <-textCh:
		case act.msg = <-clipCh:
		}
		act.dispatch()
	}
}
