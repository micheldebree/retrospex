# resize and apply b/w threshold
%.hires.1.png: %.src.png
	magick "$<" \
		-normalize \
		-threshold 50% \
		-resize 320x200^ \
		-extent 320x200 \
		-gravity center \
		"$@"

%.hires.bin: %.hires.1.png
	retrospex -f bin -m hires -o "$@" "$<"

%.hires.png: %.hires.1.png
	retrospex -m hires -o "$@" "$<"

%.hires.prg: %.hires.png
	png2prg -mode hires -display -o "$@" "$<"
