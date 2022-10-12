hook -group kakoune-pluggo global RegisterModified '"' %{
    evaluate-commands %sh{
        {{.BinPath}}/kakoune-pluggo-command cmd.put.clipboard "$kak_main_reg_dquote"
    }
}
