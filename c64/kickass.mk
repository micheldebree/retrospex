KICKASS=/Users/michel/Commodore64/Dev/KickAssembler/KickAss.jar

%.prg: %.asm $(KICKASS)
	java -jar $(KICKASS) \
		-debugdump \
		-bytedumpfile "$@.txt" \
		-symbolfile \
		-vicesymbols \
		-libdir ./lib \
		-o "$@" \
		"$<"

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
