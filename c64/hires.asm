#import "vic.asm"

.const ARTSTUDIO_TEMPLATE = "C64FILE, Bitmap=0, ScreenRam = 8000"
.const PICTURE            = LoadBinary("../scripts/paintface.hires.bin", ARTSTUDIO_TEMPLATE)
.const SCREENRAM          = $0400
.const BITMAP             = $2000

.var d011 = vic.val.d011 | vic.val.BITMAP_MODE
.var d016 = vic.val.d016
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
    ldx #0
    stx $d020
    stx $d021

loop:
    .for (var i = 0; i < 4; i++) {
        lda i * $100 + screenram_data,x
        sta i * $100 + SCREENRAM,x
    }
    inx
    bne loop
    jmp *

* = * "Screen RAM data"

screenram_data:
    .fill PICTURE.getScreenRamSize(), PICTURE.getScreenRam(i)

* = BITMAP "Bitmap data"

bitmap:
    .fill PICTURE.getBitmapSize(), PICTURE.getBitmap(i)
