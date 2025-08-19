.PHONY: all
koala.prg: koala.asm pic.bin

pic.bin: ../scripts/paintface.src.png.tmp1.png
	retrospex -f bin -o $@ $<

.PHONY: clean
clean: cleankickass
	rm *.bin

include kickass.mk
