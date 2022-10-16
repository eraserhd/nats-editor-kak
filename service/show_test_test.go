package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sends_current_session_name(t *testing.T) {
	sess := run(t, "cmd.show.data.text").OpenCommand().Session
	assert.Equal(t, "this_session", sess)
}
