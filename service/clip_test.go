package service

import (
	"testing"

	"github.com/nats-io/nats.go"
)

func Test_Updates_dquote_register_when_clip_changed(t *testing.T) {
	var receivedCmd *ChangeDquoteRegister
	act := clipChangedAction{
		msg: nats.NewMsg("event.changed.clipboard"),
		runKakouneScript: func(cmd ChangeDquoteRegister) {
			receivedCmd = &cmd
		},
	}
	act.msg.Data = []byte("foo\n")
	act.Execute()
	if receivedCmd == nil {
		t.Fatal("expected to recieve a command to change the dquote register, but did not")
	}
}
