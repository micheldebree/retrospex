#!/bin/bash
set -x
# crop and scale to 160x200 pixels
magick "$1" \
    -normalize \
    -resize 320x200^ \
    -extent 320x200 \
    -resize 160x200\! \
    -gravity center \
    "$1.0.png"

# apply koala specs
./retrospex -o "$1.1.png" "$1.0.png"

# scale back to right proportions
magick "$1.1.png" \
    -sample 320x200\! \
    "$1.2.png"

# make it into a prg
png2prg -display "$1.2.png"

