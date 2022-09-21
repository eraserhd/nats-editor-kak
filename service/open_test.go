package service

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type openOption = func(msg *nats.Msg)

func session(name string) openOption {
	return func(msg *nats.Msg) {
		msg.Header["Session"] = []string{name}
	}
}

func open(t *testing.T, opts ...openOption) OpenCmd {
	s, err := New()
	require.NoError(t, err)
	msg := &nats.Msg{
		Subject: "editor.open",
		Header: map[string][]string{
			"Session": {"editorsession"},
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return s.OpenCommand(msg)
}

func Test_Open_uses_editor_session_when_sent(t *testing.T) {
	assert.Equal(t, "foo", open(t, session("foo")).Session)
}
