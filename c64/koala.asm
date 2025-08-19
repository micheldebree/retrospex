#import "lib/vic.asm"

.file [name="%o.prg", segments="basic,code,pic"]

.segmentdef code        [startAfter = "basic"]
.segmentdef pic         [start = $2000]

.segment basic          [start = *]

BasicUpstart2(main)

.segment code

main:
    lda #$18
    sta $d018
    lda #$d8
    sta $d016
    lda #$3b
    sta $d011
    lda #0
    sta $d020
    lda backgroundcolor
    sta $d021
    ldx #0
loop:
    .for (var i = 0; i < 4; i++) {
        lda i * $100 + screenram,x
        sta i * $100 + $0400,x
        lda i * $100 + colorram,x
        sta i * $100 + $d800,x
    }
    inx
    bne loop
    jmp *

    rts

.segment pic

.label bitmap = *
.label screenram = * + 8000
.label colorram = screenram + 1000
.label backgroundcolor = colorram + 1000

.import binary "pic.bin"
