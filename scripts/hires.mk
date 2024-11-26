# USAGE:
# - name the source image <name>.src.png
# - use <name>.prg as a dependency in other targets

MODE=hires

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
%.$(MODE).png: %.1.png
	retrospex -m $(MODE) -o "$@" "$<"

# convert to prg
%.prg: %.$(MODE).png
	png2prg -mode $(MODE) -display -o "$@" "$<"
