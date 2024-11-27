# resize source to koala proportions
%.koala.1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 320x200^ \
		-extent 320x200 \
		-gravity center \
		-resize 160x200\! \
		"$@"

# apply c64 specs
%.koala.2.png: %.koala.1.png
	retrospex -m koala -o "$@" "$<"

# rescale back to double-width pixel proportions
%.koala.png: %.koala.2.png
	magick "$<" -sample 320x200\! "$@"

# convert to prg
%.koala.prg: %.koala.png
	png2prg --mode koala -display -o "$@" "$<"
