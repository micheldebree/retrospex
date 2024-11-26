# USAGE:
# - name the source image <name>.src.png
# - use <name>.prg as a dependency in other targets

MODE=mcsprites

# resize to 8x12 pixels and apply threshold
%.1.png: %.src.png
	magick "$<" \
		-normalize \
		-resize 192x \
		-resize 50%x100% \
		"$@"

# apply c64 specs
%.2.png: %.1.png retrospex
	./retrospex -m $(MODE) -o "$@" "$<"

# rescale back to double-width pixel proportions
%.$(MODE).png: %.2.png
	magick "$<" -sample 200%x100% "$@"

# convert to prg
%.prg: %.$(MODE).png
	png2prg -mode $(MODE) -display -o "$@" "$<"
