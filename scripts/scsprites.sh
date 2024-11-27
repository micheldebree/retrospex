#!/bin/bash

MODE=scsprites

# resize to 8x24 pixels and apply threshold
magick "$1" \
    -normalize \
    -resize 192x \
    -threshold 50% \
    "$1.tmp1.png"

# apply c64 specs
retrospex -m ${MODE} -o "$1.${MODE}.png" "$1.tmp1.png"

# convert to prg
png2prg -mode ${MODE} -display -o "$1.${MODE}.prg" "$1.${MODE}.png"
