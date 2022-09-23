package service

import (
	"fmt"
	"net/url"

	"github.com/nats-io/nats.go"
)

type OpenCmd struct {
	Session string
	Script  string
}

func quote(s string) string {
	result := "'"
	for _, ch := range s {
		if ch == '\'' {
			result += "'"
		}
		result += string(ch)
	}
	return result + "'"
}

func (s *Service) OpenCommand(msg *nats.Msg) OpenCmd {
	u, _ := url.Parse(string(msg.Data))
	return OpenCmd{
		Session: msg.Header.Get("Session"),
		Script:  fmt.Sprintf("edit -existing %s", quote(u.Path)),
	}
}
