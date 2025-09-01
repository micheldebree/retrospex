.PHONY: all
all: koala.prg hires.prg

koala.prg: koala.asm ../scripts/paintface.koala.bin ./lib/vic.asm

hires.prg: hires.asm ../scripts/paintface.hires.bin ./lib/vic.asm

.PHONY: clean
clean: cleankickass
	rm *.bin

include kickass.mk
include ../scripts/koala.mk
include ../scripts/hires.mk
