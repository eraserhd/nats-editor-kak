package service

import (
	"bytes"
	"log"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type SetDquoteRegister struct {
	Value string
}

var setDquoteTempl = template.Must(template.New("script").Parse(`
  define-command -override -hidden -params 1 kakoune-pluggo-set-dquote %{
    evaluate-commands %sh{
      printf 'kp entered, arg = %s\n' "$1" >&2
      printf 'kp ......., var = %s\n' "$kak_main_reg_dquote" >&2
      if [ "$1" = "$kak_main_reg_dquote" ]; then
        printf 'kp exit\n' >&2
        exit 0
      fi
      printf "set-register dquote '"
      printf %s "$1" |sed -e "s/'/''/g"
      printf "'\n"
    }
  }
  kakoune-pluggo-set-dquote {{.Value}}
`))

func (s *SetDquoteRegister) String() string {
	buf := &bytes.Buffer{}
	_ = setDquoteTempl.Execute(buf, s)
	return buf.String()
}

func executeClipChanged(a *action) {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(kakoune.Command{
		Session: a.kakouneSession,
		Script: &SetDquoteRegister{
			Value: kakoune.Quote(string(a.msg.Data)),
		},
	})
}
