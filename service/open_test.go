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
		Header:  map[string][]string{},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return s.OpenCommand(msg)
}

func Test_Defaults_session_to_kakoune(t *testing.T) {
	sess := open(t).Session
	assert.Equal(t, "kakoune", sess)
}

func Test_Open_uses_editor_session_when_sent(t *testing.T) {
	sess := open(t, header("Session", "foo")).Session
	assert.Equal(t, "foo", sess)
}

func Test_Defaults_client_to_jumpclient_option(t *testing.T) {
	client := open(t).Script.Client
	assert.Equal(t, "%opt{jumpclient}", client)
}

func Test_Allows_override_of_client_and_quotes_it(t *testing.T) {
	client := open(t, header("Window", "slime")).Script.Client
	assert.Equal(t, client, "'slime'")
}

func Test_Opens_file_URL(t *testing.T) {
	t.Run("without apostrophes", func(t *testing.T) {
		filename := open(t, data("file:///foo/bar.txt")).Script.QuotedFilename
		assert.Equal(t, filename, "'/foo/bar.txt'")
	})
	t.Run("quotes apostrophes in the filename", func(t *testing.T) {
		filename := open(t, data("file:///foo/b'ar.txt")).Script.QuotedFilename
		assert.Contains(t, filename, "'/foo/b''ar.txt'")
	})
}

func Test_Sets_editor_position(t *testing.T) {
	t.Run("defaults to line 1, column 1", func(t *testing.T) {
		script := open(t).Script
		assert.Equal(t, script.Selection, Selection{
			Start: Position{1, 1},
			End:   Position{1, 1},
		})
		assert.Equal(t, script.FixupKeys, "''")
	})
	t.Run("sets line number when given in URL", func(t *testing.T) {
		script := open(t, data("file:///foo/bar.txt#line=42")).Script
		assert.Equal(t, script.Selection, Selection{
			Start: Position{43, 1},
			End:   Position{43, 1},
		})
		assert.Equal(t, script.FixupKeys, "''")
	})
	t.Run("set line range when given in URL", func(t *testing.T) {
		script := open(t, data("file:///foo/bar.txt#line=42,47")).Script
		assert.Equal(t, script.Selection, Selection{
			Start: Position{43, 1},
			End:   Position{46, 1},
		})
		assert.Equal(t, script.FixupKeys, "'x'")
	})
}
