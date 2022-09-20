package service

import (
	"log"

	"github.com/nats-io/nats.go"
)

type Service struct {
}

func New() (*Service, error) {
	return &Service{}, nil
}

type OpenCmd struct {
	Session string
	Script  string
}

func (s *Service) OpenCommand(msg *nats.Msg) OpenCmd {
	return OpenCmd{}
}

func (s *Service) Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	ch := make(chan *nats.Msg, 32)
	sub, err := nc.ChanSubscribe("editor.open", ch)
	if err != nil {
		return err
	}
	defer sub.Drain()

	for {
		msg := <-ch
		log.Printf("recieved %q", string(msg.Data))

		_ = s.OpenCommand(msg)

		if err := msg.Respond([]byte("ok")); err != nil {
			log.Printf("error replying ok: %v", err)
			continue
		}
	}
}
