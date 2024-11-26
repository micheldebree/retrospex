#!/bin/bash
MODE=koala

# resize source to koala proportions
magick "$1" \
    -normalize \
    -resize 320x200^ \
    -extent 320x200 \
    -gravity center \
    -resize 160x200\! \
    "$1.tmp1.png"

# apply c64 specs
retrospex -m ${MODE} -o "$1.tmp2.png" "$1.tmp1.png"

# rescale back to double-width pixel proportions
magick "$1.tmp2.png" -sample 320x200\! "$1.${MODE}.png"

# convert to prg
png2prg --mode ${MODE} -display -o "$1.prg" "$1.${MODE}.png"

