#!/bin/sh

MODE=hires

# resize and apply b/w threshold
magick "$1" \
    -normalize \
    -threshold 50% \
    -resize 320x200^ \
    -extent 320x200 \
    -gravity center \
    "$1.tmp1.png"

# apply c64 specs
retrospex -m ${MODE} -o "$1.${MODE}.png" "$1.tmp1.png"

# convert to prg
png2prg -mode ${MODE} -display -o "$1.${MODE}.prg" "$1.${MODE}.png"
