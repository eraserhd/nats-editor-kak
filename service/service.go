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

type msgAction struct {
	msg *nats.Msg
}

func (a *msgAction) Execute() {
	log.Printf("recieved %q", string(a.msg.Data))
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
	if _, err := in.Write([]byte(o.Script.String())); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := in.Close(); err != nil {
		return fmt.Errorf("closing pipe: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error responding: %w", err)
	}
	return nil
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
			action := msgAction{msg}
			action.Execute()

			open := OpenCommand(msg)
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
		case <-clipCh:
			log.Printf("clipboard changed")
		}
	}
}
