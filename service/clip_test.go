package service

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/plugbench/kakoune-pluggo/kakoune"
)

func Test_Updates_dquote_register_when_clip_changed(t *testing.T) {
	var receivedCmd *kakoune.Command
	act := clipChangedAction{
		msg: nats.NewMsg("event.changed.clipboard"),
		runKakouneScript: func(cmd kakoune.Command) error {
			receivedCmd = &cmd
			return nil
		},
	}
	act.msg.Data = []byte("foo\n")
	act.Execute()
	if receivedCmd == nil {
		t.Fatal("expected to recieve a command to change the dquote register, but did not")
	}
	_, ok := receivedCmd.Script.(*SetDquoteRegister)
	if !ok {
		t.Fatalf("expected kakoune script to be *SetDquoteRegister, but was %T", receivedCmd.Script)
	}
}
