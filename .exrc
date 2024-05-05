set nolist
set fdm=syntax
set fdl=999
set relativenumber
" nmap <F5> :!go run *.go<CR>
" nmap <F5> :wa<cr>:vsplit <bar> term make<cr>
nmap <F5> :wa<cr>:make<cr>zz
nmap <S-F5> :wa<cr>:vsplit<bar>term make<cr>
" packadd vim-go
