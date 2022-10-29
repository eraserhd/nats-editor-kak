package service

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/plugbench/kakoune-pluggo/kakoune"
)

type showText struct {
	Client string
	Text   string
}

var showTextTempl = template.Must(template.New("script").Parse(`
  evaluate-commands -save-regs t -try-client {{.Client}} %{
    try %{
      evaluate-commands %sh{
        have_show=false
        next_n=0
        eval set -- "$kak_quoted_buflist"
        for buf in "$@"; do
          case "$buf" in
            "*show*")
              have_show=true
              ;;
            "*show-"*"*")
              n_part=${buf%"*"}
              n_part=${n_part#"*"show-}
              next_n=$(( n_part >= next_n ? n_part + 1 : next_n ))
              ;;
          esac
        done
        if [ "$have_show" = false ]; then
          printf 'edit -scratch *show*\n'
        else
          printf 'edit -scratch *show-%d*\n' "$next_n"
        fi
      }
      set-register t {{.Text}}
      execute-keys '%"tR'
      try focus
    } catch %{
      echo -markup "{Error}%val{error}"
      echo -debug "%val{error}"
    }
  }
`))

func (s *showText) String() string {
	buf := &bytes.Buffer{}
	showTextTempl.Execute(buf, s)
	return buf.String()
}

func executeShowText(act *action) {
	act.log("info", fmt.Sprintf("received text to show: %q", string(act.msg.Data)))
	cmd := kakoune.Command{
		Session: act.kakouneSession,
		Script: &showText{
			Client: "%opt{jumpclient}",
			Text:   kakoune.Quote(string(act.msg.Data)),
		},
	}
	act.runKakouneScript(cmd)
}
