package service

import (
	"fmt"
	"log"
	"os/exec"

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

func (o OpenCmd) Run(msg *nats.Msg) error {
	cmd := exec.Command("kak", "-p", o.Session)
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting kak: %w", err)
	}
	if _, err := in.Write([]byte(o.Script)); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error responding: %w", err)
	}
	return nil
}

func (s *Service) OpenCommand(msg *nats.Msg) OpenCmd {
	return OpenCmd{Session: msg.Header.Get("Session")}
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

		open := s.OpenCommand(msg)
		if err := open.Run(msg); err != nil {
        		log.Print(err)
			if err := msg.Respond([]byte(fmt.Sprintf("ERROR: %s", err.Error()))); err != nil {
				log.Printf("error responding: %v", err)
			}
			continue
		}
		if err := msg.Respond([]byte("ok")); err != nil {
			log.Printf("error replying ok: %v", err)
			continue
		}
	}
}
