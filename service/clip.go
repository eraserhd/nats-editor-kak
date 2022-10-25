package service

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type SetDquoteRegister struct {
	BinPath string
	Value   string
}

var setDquoteTempl = template.Must(template.New("script").Parse(`
  define-command -override -hidden -params 1 kakoune-pluggo-set-dquote %{
    evaluate-commands %sh{
      {{.BinPath}}/kakoune-pluggo-event 'event.logged.kakoune-pluggo.debug' "setting from '$kak_main_reg_dquote' to '$1'" 2>/dev/null
      if [ "$1" = "$kak_main_reg_dquote" ]; then
        {{.BinPath}}/kakoune-pluggo-event 'event.logged.kakoune-pluggo.debug' "skipping update" 2>/dev/null
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
