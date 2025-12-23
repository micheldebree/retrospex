.PHONY: all
all: koala.prg hires.prg mixedhires.prg

koala.prg: koala.asm paintface.koala.bin ./lib/vic.asm

hires.prg: hires.asm paintface.hires.bin ./lib/vic.asm

mixedhires.prg: mixedhires.asm paintface.mixedhires.bin ./lib/vic.asm

.PHONY: clean
clean: cleankickass
	rm *.bin || true

include kickass.mk
include ../scripts/koala.mk
include ../scripts/hires.mk
