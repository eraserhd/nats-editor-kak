package service

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type openOption = func(msg *nats.Msg)

func header(name, value string) openOption {
	return func(msg *nats.Msg) {
		msg.Header[name] = []string{value}
	}
}

func data(data string) openOption {
	return func(msg *nats.Msg) {
		msg.Data = []byte(data)
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
	assert.Equal(t, "foo", open(t, header("Session", "foo")).Session)
}

func Test_Defaults_client_to_jumpclient_option(t *testing.T) {
	client := open(t).Script.Client
	assert.Equal(t, client, "%opt{jumpclient}")
}

func Test_Allows_override_of_client_and_quotes_it(t *testing.T) {
	client := open(t, header("Window", "slime")).Script.Client
	assert.Equal(t, client, "'slime'")
}

func Test_Opens_file_URL(t *testing.T) {
	assert.Equal(t,
		open(t, data("file:///foo/bar.txt")).Script.QuotedFilename,
		"'/foo/bar.txt'",
	)
	t.Run("quotes apostrophes in the filename", func(t *testing.T) {
		assert.Contains(t,
			open(t, data("file:///foo/b'ar.txt")).Script.QuotedFilename,
			"'/foo/b''ar.txt'",
		)
	})
}
