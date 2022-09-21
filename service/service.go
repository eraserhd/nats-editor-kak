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

		cmd := exec.Command("kak", "-p", open.Session)
		in, err := cmd.StdinPipe()
		if err != nil {
			log.Printf("error creating pipe: %v", err)
			if err := msg.Respond([]byte(fmt.Sprintf("ERROR: %s", err.Error()))); err != nil {
				log.Printf("error responding: %v", err)
			}
			continue
		}
		if err := cmd.Start(); err != nil {
			log.Printf("error starting kak: %v", err)
			if err := msg.Respond([]byte(fmt.Sprintf("ERROR: %s", err.Error()))); err != nil {
				log.Printf("error responding: %v", err)
			}
			continue
		}
		if _, err := in.Write([]byte(open.Script)); err != nil {
			log.Printf("error writing script: %v", err)
			if err := msg.Respond([]byte(fmt.Sprintf("ERROR: %s", err.Error()))); err != nil {
				log.Printf("error responding: %v", err)
			}
			continue
		}
		if err := cmd.Wait(); err != nil {
			log.Printf("error waiting for kak: %v", err)
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
