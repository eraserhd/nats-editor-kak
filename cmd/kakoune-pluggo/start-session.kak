declare-option -hidden str pluggo_last_yank_client %val{client}

hook -group kakoune-pluggo global RegisterModified '"' %{
    set-option global pluggo_last_yank_client %val{client}
    evaluate-commands %sh{
        {{.BinPath}}/kakoune-pluggo command cmd.put.clipboard "$kak_main_reg_dquote"
    }
}

evaluate-commands %sh{
    {{.BinPath}}/kakoune-pluggo-daemon "$kak_session" </dev/null >/dev/null 2>&1 &
    daemon_pid=$!
    printf 'declare-option -hidden str pluggo_daemon_pid "%s"\n' "$daemon_pid"
    {{.BinPath}}/kakoune-pluggo-event "event.logged.kakoune-pluggo.info" "pid for session $kak_session is $daemon_pid"
}

hook -group kakoune-pluggo global KakEnd .* %{
    nop %sh{ kill -HUP "$kak_opt_pluggo_daemon_pid" >/dev/null 2>&1 }
}

define-command -hidden -params 1 kakoune-pluggo-set-dquote %{
    evaluate-commands -try-client %opt{pluggo_last_yank_client} %{
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
}
