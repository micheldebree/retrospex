KICKASS=/Users/michel/Commodore64/Dev/KickAssembler/KickAss.jar
VICE=/Applications/vice-arm64-gtk3-3.9/bin/x64sc
DEBUGGER=/Applications/Retro\ Debugger.app/Contents/MacOS/Retro\ Debugger
EXOMIZER=/usr/local/bin/exomizer

%.prg: %.asm $(KICKASS)
	java -jar $(KICKASS) \
		-debugdump \
		-bytedumpfile "$@.txt" \
		-symbolfile \
		-vicesymbols \
		-libdir ./lib \
		-log log.txt \
		-o "$@" \
		"$<"

.PRECIOUS: %.exe.prg
%.exe.prg: %.prg
	exomizer sfx basic "$<" -o "$@"

.PHONY: run
run: main.prg
	$(VICE) \
		-autostartprgmode 1 \
		-autostart-warp \
		-moncommands main.vs \
		-VICIIborders debug \
		"$<"

.PHONY: debug
debug: main.prg
	$(DEBUGGER) \
		-wait 2500 \
		-autojmp \
		-layout 9 \
		-debuginfo main.dbg \
		 -prg "$<"

.PHONY: cleankickass
cleankickass:
	rm -f *.prg
	rm -f *.exe.prg
	rm -f *.sym
	rm -f *.vs
	rm -f *.dbg
	rm -f *.d64
	rm -f *.out
	rm -f *.map
	rm -f *.labels
	rm -f *.txt
