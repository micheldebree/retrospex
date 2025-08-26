#import "vic.asm"

.const PICTURE      = LoadBinary("pic.bin", BF_KOALA)
.const SCREENRAM    = $0400
.const BITMAP       = $2000

.var d011 = vic.val.d011 | vic.val.BITMAP_MODE
.var d016 = vic.val.d016 | vic.val.MULTICOLOR_MODE
.var d018 = @vic_bitmap(bitmap) | @vic_screenram(SCREENRAM)
.var dd00 = vic.val.BANK_0

BasicUpstart2(main)

* = * "Code"

main:

    lda #dd00
    sta $dd00
    lda #d011
    sta $d011
    lda #d016
    sta $d016
    lda #d018
    sta $d018
    lda #PICTURE.getBackgroundColor()
    sta $d020
    sta $d021
    ldx #0

loop:
    .for (var i = 0; i < 4; i++) {
        lda i * $100 + screenram_data,x
        sta i * $100 + SCREENRAM,x
        lda i * $100 + colorram_data,x
        sta i * $100 + vic.COLOR_RAM,x
    }
    inx
    bne loop
    jmp *

* = * "Color RAM data"

colorram_data:
    .fill PICTURE.getColorRamSize(), PICTURE.getColorRam(i)

* = * "Screen RAM data"

screenram_data:
    .fill PICTURE.getScreenRamSize(), PICTURE.getScreenRam(i)

* = BITMAP "Bitmap"

bitmap:
    .fill PICTURE.getBitmapSize(), PICTURE.getBitmap(i)
