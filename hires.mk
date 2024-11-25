# USAGE:
# - name the source image <name>.src.png
# - use <name>.prg as a dependency in other targets

# resize and apply b/w threshold
%.1.png: %.src.png
	magick "$<" \
		-normalize \
		-threshold 50% \
		-resize 320x200^ \
		-extent 320x200 \
		-gravity center \
		"$@"

# apply c64 specs
%.2.png: %.1.png ./retrospex
	./retrospex -m "hires" -o "$@" "$<"

# rescale back to double-width pixel proportions
.PRECIOUS: %.hires.png
%.hires.png: %.2.png
	magick "$<" -sample 320x200\! "$@"

# convert to prg
%.prg: %.hires.png
	png2prg -display -o "$@" "$<"

