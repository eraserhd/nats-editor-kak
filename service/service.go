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
	log              func(level, text string)
}

var channelHandlers = []struct {
	subject string
	handler func(a *action)
}{
	{"cmd.show.url.file", executeShowFileURL},
	{"cmd.show.data.text", executeShowText},
	{"event.changed.clipboard", executeClipChanged},
}

func (a *action) dispatch() {
	for _, h := range channelHandlers {
		if a.msg.Subject == h.subject {
			h.handler(a)
			return
		}
	}
	log.Fatalf("do not know how to handle %s", a.msg.Subject)
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

	chs := make([]chan *nats.Msg, len(channelHandlers))
	for i, h := range channelHandlers {
		chs[i] = make(chan *nats.Msg, 32)
		sub, err := nc.ChanSubscribe(h.subject, chs[i])
		if err != nil {
			return err
		}
		defer sub.Drain()
	}

	for {
		act := action{
			kakouneSession: s.kakouneSession,
			publish: func(msg *nats.Msg) error {
				return nc.PublishMsg(msg)
			},
			runKakouneScript: kakoune.Run,
			log: func(level, text string) {
				msg := nats.NewMsg("event.logged.kakoune-pluggo." + level)
				msg.Data = []byte(text)
				nc.PublishMsg(msg)
			},
		}
		select {
		case act.msg = <-chs[0]:
		case act.msg = <-chs[1]:
		case act.msg = <-chs[2]:
		}
		act.dispatch()
	}
}
