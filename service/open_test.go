package service

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Open_uses_editor_session_when_sent(t *testing.T) {
	s, err := New()
	require.NoError(t, err)
	cmd := s.OpenCommand(&nats.Msg{
        	Subject: "editor.open",
        	Header: map[string][]string{
                	"Session": {"foo"},
        	},
	})
	assert.Equal(t, "foo", cmd.Session)
}
