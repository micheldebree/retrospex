#!/bin/bash

MODE=mcsprites

# resize to 96 (8x12) pixels width
magick "$1" \
    -normalize \
    -resize 192x \
    -resize 50%x100% \
    "$1.tmp1.png"

# apply c64 specs
retrospex -m ${MODE} -o "$1.tmp2.png" "$1.tmp1.png"

# rescale back to double-width pixel proportions
magick "$1.tmp2.png" -sample 200%x100% "$1.${MODE}.png"

# convert to prg
png2prg -mode ${MODE} -display -o "$1.${MODE}.prg" "$1.${MODE}.png"
