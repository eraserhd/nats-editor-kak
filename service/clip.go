package service

import (
	"bytes"
	"log"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type SetDquoteRegister struct {
	Value string
}

var setDquoteTempl = template.Must(template.New("script").Parse(`
  define-command -override -hidden -params 1 kakoune-pluggo-set-dquote %{
    evaluate-commands %sh{
      if [ "$1" = "$kak_main_reg_dquote" ]; then
        exit 0
      fi
      printf "set-register dquote '"
      printf %s "$kak_main_reg_dquote" |sed -e "s/'/''/g"
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

type clipChangedAction struct {
	msg              *nats.Msg
	runKakouneScript func(cmd kakoune.Command) error
}

func (a *clipChangedAction) Execute() {
	log.Printf("recieved clipboard changed event")
	a.runKakouneScript(kakoune.Command{
		Session: "kakoune",
		Script: &SetDquoteRegister{
			Value: kakoune.Quote(string(a.msg.Data)),
		},
	})
}
