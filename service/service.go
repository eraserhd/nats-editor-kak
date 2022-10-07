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
	msg            *nats.Msg
	publish        func(msg *nats.Msg) error
	runOpenCommand func(cmd OpenCommand) error
}

func (a *msgAction) Execute() {
	log.Printf("recieved %q", string(a.msg.Data))

	cmd := openCommand(a.msg)
	if err := a.runOpenCommand(cmd); err != nil {
		log.Print(err)
		reply := nats.NewMsg(a.msg.Reply)
		reply.Data = []byte(fmt.Sprintf("ERROR: %s", err.Error()))
		if err := a.publish(reply); err != nil {
			log.Printf("error responding: %v", err)
		}
		return
	}

	reply := nats.NewMsg(a.msg.Reply)
	reply.Data = []byte("ok")
	if err := a.publish(reply); err != nil {
		log.Printf("error replying ok: %v", err)
		return
	}
}

func runOpenCommand(o OpenCommand) error {
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
			action := msgAction{
				msg: msg,
				publish: func(msg *nats.Msg) error {
					return nc.PublishMsg(msg)
				},
				runOpenCommand: runOpenCommand,
			}
			action.Execute()
		case <-clipCh:
			log.Printf("clipboard changed")
		}
	}
}
