hook -group kakoune-pluggo global RegisterModified '"' %{
    evaluate-commands %sh{
        {{.BinPath}}/kakoune-pluggo-command cmd.put.clipboard "$kak_main_reg_dquote"
    }
}

evaluate-commands %sh{
    {{.BinPath}}/kakoune-pluggo-daemon "$kak_session" </dev/null >/dev/null 2>&1 &
    daemon_pid=$!
    printf 'declare-option -hidden str pluggo_daemon_pid "%s"\n' "$daemon_pid"
}

hook -group kakoune-pluggo global KakEnd .* %{
    nop %sh{ kill -HUP "$kak_opt_pluggo_daemon_pid" >/dev/null 2>&1 }
}
