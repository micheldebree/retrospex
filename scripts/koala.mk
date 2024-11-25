# USAGE:
# - name the source image <name>.src.png
# - use <name>.prg as a dependency in other targets

# resize source to koala proportions
%.1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 320x200^ \
		-extent 320x200 \
		-gravity center \
		-resize 160x200\! \
		"$@"

# apply c64 specs
%.2.png: %.1.png ./retrospex
	./retrospex -o "$@" "$<"

# rescale back to double-width pixel proportions
.PRECIOUS: %.koala.png
%.koala.png: %.2.png
	magick "$<" -sample 320x200\! "$@"

# convert to prg
%.prg: %.koala.png
	png2prg -display -o "$@" "$<"

