package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (r runResult) showTextScript() *showText {
	script, ok := r.kakouneCommand().Script.(*showText)
	if !ok {
		r.t.Fatalf("wanted script to be of type *showText, got %T", r.kakouneCommand().Script)
	}
	return script
}

func Test_sends_current_session_name(t *testing.T) {
	sess := run(t, "cmd.show.data.text").kakouneCommand().Session
	assert.Equal(t, "this_session", sess)
}

func Test_defaults_client_to_jumpclient(t *testing.T) {
	client := run(t, "cmd.show.data.text").showTextScript().Client
	assert.Equal(t, "%opt{jumpclient}", client)
}

func Test_sends_quoted_text(t *testing.T) {
	text := run(t, "cmd.show.data.text", data("hell'o world")).showTextScript().Text
	assert.Equal(t, "'hell''o world'", text)
}

func Test_sends_ok_reply(t *testing.T) {
        msg := run(t, "cmd.show.data.text", data("hell'o world")).Reply()
        assert.Equal(t, "ok", string(msg.Data))
}
