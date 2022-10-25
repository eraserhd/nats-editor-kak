package service

import (
	"bytes"
	"fmt"
	"text/template"

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
	a.log("info", fmt.Sprintf("clipboard changed: %q", string(a.msg.Data)))
	a.runKakouneScript(kakoune.Command{
		Session: a.kakouneSession,
		Script: &SetDquoteRegister{
			Value: kakoune.Quote(string(a.msg.Data)),
		},
	})
}
