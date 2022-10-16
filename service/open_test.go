package service

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"github.com/plugbench/kakoune-pluggo/kakoune"
	"github.com/plugbench/kakoune-pluggo/service/fragment"
)

type runResult struct {
	t                    *testing.T
	msg                  *nats.Msg
	scriptExecutionError error

	executedScripts   []kakoune.Command
	publishedMessages []*nats.Msg
}

func (result runResult) OpenCommand() kakoune.Command {
	if len(result.executedScripts) != 1 {
		result.t.Fatalf("expected 1 script to be executed, but got %d", len(result.executedScripts))
	}
	return result.executedScripts[0]
}

func (result runResult) OpenScript() OpenFile {
	sc := result.OpenCommand().Script
	if sc, ok := sc.(*OpenFile); ok {
		return *sc
	}
	result.t.Fatal("script was of the wrong type")
	return OpenFile{}
}

func (result runResult) Reply() *nats.Msg {
	var replies []*nats.Msg
	for _, msg := range result.publishedMessages {
		if msg.Subject == "_INBOX.Reply" {
			replies = append(replies, msg)
		}
	}
	if len(replies) != 1 {
		result.t.Fatalf("expected 1 reply, but got %d", len(replies))
	}
	return replies[0]
}

type runOption = func(result *runResult)

func header(name, value string) runOption {
	return func(result *runResult) { result.msg.Header[name] = []string{value} }
}

func data(data string) runOption {
	return func(result *runResult) { result.msg.Data = []byte(data) }
}

func scriptExecutionError(text string) runOption {
	return func(result *runResult) {
		result.scriptExecutionError = errors.New(text)
	}
}

func run(t *testing.T, opts ...runOption) runResult {
	result := runResult{
		t: t,
		msg: &nats.Msg{
			Reply:   "_INBOX.Reply",
			Subject: "editor.open",
			Header:  map[string][]string{},
		},
	}
	for _, opt := range opts {
		opt(&result)
	}
	act := action{
		kakouneSession: "this_session",
		msg:            result.msg,
		publish: func(msg *nats.Msg) error {
			result.publishedMessages = append(result.publishedMessages, msg)
			return nil
		},
		runKakouneScript: func(cmd kakoune.Command) error {
			if result.scriptExecutionError != nil {
				return result.scriptExecutionError
			}
			result.executedScripts = append(result.executedScripts, cmd)
			return nil
		},
		execute: executeShowFileURL,
	}
	act.execute(&act)
	return result
}

func Test_Sends_current_session_name(t *testing.T) {
	sess := run(t).OpenCommand().Session
	assert.Equal(t, "this_session", sess)
}

func Test_Defaults_client_to_jumpclient_option(t *testing.T) {
	client := run(t).OpenScript().Client
	assert.Equal(t, "%opt{jumpclient}", client)
}

func Test_Allows_override_of_client_and_quotes_it(t *testing.T) {
	client := run(t, header("Window", "slime")).OpenScript().Client
	assert.Equal(t, client, "'slime'")
}

func Test_Opens_file_URL(t *testing.T) {
	t.Run("without apostrophes", func(t *testing.T) {
		filename := run(t, data("file:///foo/bar.txt")).OpenScript().QuotedFilename
		assert.Equal(t, filename, "'/foo/bar.txt'")
	})
	t.Run("quotes apostrophes in the filename", func(t *testing.T) {
		filename := run(t, data("file:///foo/b'ar.txt")).OpenScript().QuotedFilename
		assert.Contains(t, filename, "'/foo/b''ar.txt'")
	})
}

func Test_Sets_editor_position(t *testing.T) {
	t.Run("defaults to line 1, column 1", func(t *testing.T) {
		script := run(t).OpenScript()
		assert.Equal(t, script.Selection, fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 1, Column: 1},
			End:   fragment.LineAndColumn{Line: 1, Column: 1},
		})
		assert.Equal(t, script.FixupKeys, "''")
	})
	t.Run("sets line number when given in URL", func(t *testing.T) {
		script := run(t, data("file:///foo/bar.txt#line=42")).OpenScript()
		assert.Equal(t, script.Selection, fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 43, Column: 1},
			End:   fragment.LineAndColumn{Line: 43, Column: 1},
		})
		assert.Equal(t, script.FixupKeys, "''")
	})
	t.Run("sets line and column number when given in URL", func(t *testing.T) {
		script := run(t, data("file:///foo/bar.txt#line=42.3")).OpenScript()
		assert.Equal(t, script.Selection, fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 43, Column: 4},
			End:   fragment.LineAndColumn{Line: 43, Column: 4},
		})
		assert.Equal(t, script.FixupKeys, "''")
	})
	t.Run("set line range when given in URL", func(t *testing.T) {
		script := run(t, data("file:///foo/bar.txt#line=2,5")).OpenScript()
		assert.Equal(t, script.Selection, fragment.LineAndColumnSelection{
			Start: fragment.LineAndColumn{Line: 3, Column: 1},
			End:   fragment.LineAndColumn{Line: 5, Column: 1},
		})
		assert.Equal(t, script.FixupKeys, "'<a-L>'")
	})
}

func Test_Sends_replies(t *testing.T) {
	t.Run("ok reply when everything works", func(t *testing.T) {
		reply := run(t).Reply()
		assert.Equal(t, string(reply.Data), "ok")
	})
	t.Run("ERROR reply when the editor command fails", func(t *testing.T) {
		reply := run(t, scriptExecutionError("command failed")).Reply()
		assert.Equal(t, string(reply.Data), "ERROR: command failed")
	})
}
