let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/Documents/School/ADET/Finals/Bazaroo/server
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +1 ~/Documents/School/ADET/Finals/Bazaroo/server/server/routes.go
badd +23 ~/Documents/School/ADET/Finals/Bazaroo/server/server/api.go
badd +18 ~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go
badd +69 ~/Documents/School/ADET/Finals/Bazaroo/server/server/handlers/addresses.go
badd +58 ~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go
badd +137 ~/Documents/School/ADET/Finals/Bazaroo/server/server/utils/validation.go
argglobal
%argdel
edit ~/Documents/School/ADET/Finals/Bazaroo/server/server/utils/validation.go
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
split
1wincmd k
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd w
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
wincmd =
argglobal
balt ~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 137 - ((7 * winheight(0) + 8) / 16)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 137
normal! 035|
wincmd w
argglobal
if bufexists(fnamemodify("~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go", ":p")) | buffer ~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go | else | edit ~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go | endif
if &buftype ==# 'terminal'
  silent file ~/Documents/School/ADET/Finals/Bazaroo/server/server/db/queries.go
endif
balt ~/Documents/School/ADET/Finals/Bazaroo/server/server/handlers/addresses.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 53 - ((9 * winheight(0) + 8) / 16)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 53
normal! 0
wincmd w
argglobal
if bufexists(fnamemodify("~/Documents/School/ADET/Finals/Bazaroo/server/server/api.go", ":p")) | buffer ~/Documents/School/ADET/Finals/Bazaroo/server/server/api.go | else | edit ~/Documents/School/ADET/Finals/Bazaroo/server/server/api.go | endif
if &buftype ==# 'terminal'
  silent file ~/Documents/School/ADET/Finals/Bazaroo/server/server/api.go
endif
balt ~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 23 - ((8 * winheight(0) + 8) / 17)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 23
normal! 020|
wincmd w
argglobal
if bufexists(fnamemodify("~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go", ":p")) | buffer ~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go | else | edit ~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go | endif
if &buftype ==# 'terminal'
  silent file ~/Documents/School/ADET/Finals/Bazaroo/server/server/conndb.go
endif
balt ~/Documents/School/ADET/Finals/Bazaroo/server/server/routes.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 18 - ((8 * winheight(0) + 8) / 17)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 18
normal! 0
wincmd w
wincmd =
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
