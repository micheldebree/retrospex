# resize to 8x24 pixels and apply threshold
%.tmp1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 192x \
		-threshold 50% \
		"$@"

# apply c64 specs
%.scsprites.png: %.tmp1.png ./retrospex
	./retrospex -m scsprites -o "$@" "$<"

# convert to prg
%.scsprites.prg: %.scsprites.png
	png2prg -mode scsprites -display -o "$@" "$<"

