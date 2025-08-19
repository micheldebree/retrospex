set errorformat=%EError:\ %m,%Cat\ line\ %l\\,\ column\ %c\ in\ %f,%Z
set autoindent
set nosmartindent
set shiftwidth=4
set tabstop=4
set softtabstop=4
set expandtab
set foldlevel=0
set foldcolumn=3
set number
set relativenumber
au BufRead *.asm set ft=kickass
au FileType kickass set expandtab
au FileType kickass set list
au FileType kickass set commentstring=//%s
au FileType kickass setlocal foldmarker={,}
au FileType kickass set foldmethod=marker
au FileType make set noexpandtab
au FileType make set vartabstop=
au FileType make set shiftwidth=4
au FileType make set list
