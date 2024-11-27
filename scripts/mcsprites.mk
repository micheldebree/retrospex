# resize to 8x12 pixels and apply threshold
%.tmp1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 192x \
		-resize 50%x100% \
		"$@"

# apply c64 specs
%.tmp2.png: %.tmp1.png retrospex
	./retrospex -m mcsprites -o "$@" "$<"

# rescale back to double-width pixel proportions
%.mcsprites.png: %.tmp2.png
	magick "$<" -sample 200%x100% "$@"

# convert to prg
%.mcsprites.prg: %.mcsprites.png
	png2prg -mode mcsprites -display -o "$@" "$<"
