package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"text/template"

	"github.com/nats-io/nats.go"
)

type Script struct {
	Client         string
	QuotedFilename string
	Selection      Selection
}

type Selection struct {
	Start, End Position
}

type Position struct {
	Line, Column int
}

var templ = template.Must(template.New("script").Parse(`
  evaluate-commands -try-client {{.Client}} %{
    try %{
      edit -existing {{.QuotedFilename}}
      try focus
    } catch %{
      echo -markup "{Error}%val{error}"
      echo -debug "%val{error}"
    }
  }
`))

var fragmentRegexp = regexp.MustCompile(`^line=(\d+)`)

func (s Script) String() string {
	buf := &bytes.Buffer{}
	_ = templ.Execute(buf, s)
	return buf.String()
}

type OpenCmd struct {
	Session string
	Script  Script
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

func parseFragment(fragment string) (int, error) {
	match := fragmentRegexp.FindStringSubmatch(fragment)
	if match == nil {
		return 0, errors.New("cannot parse fragment identifier")
	}
	line, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing fragment identifer: %w", err)
	}
	return int(line), nil
}

func (s *Service) OpenCommand(msg *nats.Msg) OpenCmd {
	u, _ := url.Parse(string(msg.Data))
	result := OpenCmd{
		Session: "kakoune",
		Script: Script{
			Client:         "%opt{jumpclient}",
			QuotedFilename: quote(u.Path),
			Selection: Selection{
				Start: Position{1, 1},
				End:   Position{1, 1},
			},
		},
	}
	if line, err := parseFragment(u.Fragment); err == nil {
		result.Script.Selection = Selection{
			Start: Position{int(line) + 1, 1},
			End:   Position{int(line) + 1, 1},
		}
	}
	if s := msg.Header.Get("Session"); s != "" {
		result.Session = s
	}
	if w := msg.Header.Get("Window"); w != "" {
		result.Script.Client = quote(w)
	}
	return result
}
