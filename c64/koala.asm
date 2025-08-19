#import "lib/vic.asm"

.const SCREENRAM = $0400
.const KOALA_TEMPLATE = "Bitmap=$0000, ScreenRam=$1f40, ColorRam=$2328, BackgroundColor = $2710"
.var picture = LoadBinary("pic.bin", KOALA_TEMPLATE)

.file [name="%o.prg", segments="basic,code,data,bitmap"]

.segmentdef code        [startAfter = "basic"]
.segmentdef data        [startAfter = "code"]
.segmentdef bitmap      [start = $2000]

.segment basic          [start = *]

BasicUpstart2(main)

.segment code

main:

    @vic_memory_bank(vic.memory.BANK_0000)
    lda #@vic_layout(bitmap, SCREENRAM)
    sta $d018

    lda #@vic_control2_value(0, false, true)
    sta $d016
    lda #$3b
    sta $d011
    lda #picture.getBackgroundColor()
    sta $d020
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

.segment data

colorram:
    .fill picture.getColorRamSize(), picture.getColorRam(i)

screenram:
    .fill picture.getScreenRamSize(), picture.getScreenRam(i)

.segment bitmap
bitmap:
    .fill picture.getBitmapSize(), picture.getBitmap(i)



