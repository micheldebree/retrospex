# USAGE:
# - name the source image <name>.src.png
# - use <name>.prg as a dependency in other targets

MODE=scsprites

# resize to 8x24 pixels and apply threshold
%.1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 192x \
		-threshold 50% \
		"$@"

# apply c64 specs
%.$(MODE).png: %.1.png ./retrospex
	./retrospex -m $(MODE) -o "$@" "$<"

# convert to prg
%.prg: %.$(MODE).png
	png2prg -mode $(MODE) -display -o "$@" "$<"

