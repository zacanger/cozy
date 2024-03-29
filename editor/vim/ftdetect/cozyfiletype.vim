" don't spam the user when Vim is started in Vi compatibility mode
let s:cpo_save = &cpo
set cpo&vim

" Note: should not use augroup in ftdetect (see :help ftdetect)
au BufRead,BufNewFile *.cz setfiletype cozy
au BufRead,BufNewFile .cozy_init setfiletype cozy

" restore Vi compatibility settings
let &cpo = s:cpo_save
unlet s:cpo_save
