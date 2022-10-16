package service

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

func Test_Updates_dquote_register_when_clip_changed(t *testing.T) {
	var receivedCmd *kakoune.Command
	act := action{
		kakouneSession: "foosess",
		msg:            nats.NewMsg("event.changed.clipboard"),
		runKakouneScript: func(cmd kakoune.Command) error {
			receivedCmd = &cmd
			return nil
		},
	}
	act.msg.Data = []byte("foo\n")
	act.dispatch()
	require.NotNil(t, receivedCmd, "expected to recieve a command to change the dquote register, but did not")
	setDquote, ok := receivedCmd.Script.(*SetDquoteRegister)
	require.Truef(t, ok, "expected kakoune script to be *SetDquoteRegister, but was %T", receivedCmd.Script)
	assert.Equal(t, receivedCmd.Session, "foosess")
	assert.Equal(t, setDquote.Value, "'foo\n'")
}
