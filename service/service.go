package service

import (
	"fmt"
	"os/exec"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type Service struct {
}

func New() (*Service, error) {
	return &Service{}, nil
}

func runKakouneScript(o kakoune.Command) error {
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
			action := openAction{
				msg: msg,
				publish: func(msg *nats.Msg) error {
					return nc.PublishMsg(msg)
				},
				runKakouneScript: runKakouneScript,
			}
			action.Execute()
		case msg := <-clipCh:
			action := clipChangedAction{
				msg: msg,
				runKakouneScript: runKakouneScript,
			}
			action.Execute()
		}
	}
}
