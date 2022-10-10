hook -group kakoune-pluggo global RegisterModified '"' %{
    evaluate-commands %sh{
        {{.BinPath}}/kakoune-pluggo-command cmd.clipboard.put "$kak_main_reg_dquote"
    }
}
