.PHONY: all
all: paintface.prg

%.1.png: %.0.png
	magick "$<" \
		-normalize \
		-resize 320x200^ \
		-extent 320x200 \
		-gravity center \
		-resize 160x200\! \
		"$@"

%.2.png: %.1.png
	./retrospex -o "$@" "$<"

%.3.png: %.2.png
	magick "$<" -sample 320x200\! "$@"

%.prg: %.3.png
	png2prg -display -o "$@" "$<"



