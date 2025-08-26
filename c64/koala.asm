#import "lib/vic.asm"

.const SCREENRAM = $0400
.var picture = LoadBinary("pic.bin", BF_KOALA)

.file [name="%o.prg", segments="basic,code,data,bitmap"]

.segmentdef code        [startAfter = "basic"]
.segmentdef data        [startAfter = "code"]
.segmentdef bitmap      [start = $2000]

.segment basic          [start = *]

BasicUpstart2(main)

.segment code

main:
    .var d011 = vic.default.d011 | vic.bit.BITMAP_MODE
    .var d016 = vic.default.d016 | vic.bit.MULTICOLOR_MODE
    .var d018 = @vic_bitmap(bitmap) | @vic_screenram(SCREENRAM)

    @vic_memory_bank(vic.memory.BANK_0000)
    lda #d011
    sta $d011
    lda #d016
    sta $d016
    lda #d018
    sta $d018
    lda #picture.getBackgroundColor()
    sta $d020
    sta $d021
    ldx #0
loop:
    .for (var i = 0; i < 4; i++) {
        lda i * $100 + screenram_data,x
        sta i * $100 + SCREENRAM,x
        lda i * $100 + colorram_data,x
        sta i * $100 + vic.register.COLOR_RAM,x
    }
    inx
    bne loop
    jmp *

    rts

.segment data

colorram_data:
    .fill picture.getColorRamSize(), picture.getColorRam(i)

screenram_data:
    .fill picture.getScreenRamSize(), picture.getScreenRam(i)

.segment bitmap
bitmap:
    .fill picture.getBitmapSize(), picture.getBitmap(i)
