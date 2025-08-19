#importonce

.filenamespace vic

register: {
    .label CONTROL1         = $d011
    .label SPRITE_ENABLE    = $d015
    .label CONTROL2         = $d016
    .label MEM_LAYOUT       = $d018
    .label BORDER           = $d020
    .label BACKGROUND       = $d021
    .label MULTICOLOR_1     = $d022
    .label MULTICOLOR_2     = $d023
    .label SPRITE_COLOR     = $d027
    .label COLOR_MEM        = $d800
    .label BANK_SELECT      = $dd00
}

// some common VIC data sizes
size: {
    .label SPRITE           = $40
    .label FONT             = $800
    .label BANK             = $4000
    .label CHARS_PER_LINE   = 40 
    .label CHAR             = 8
}

memory: {
    .label BANK_0000        = %11
    .label BANK_4000        = %10
    .label BANK_8000        = %01
    .label BANK_C000        = %00
}

// select one of the four vic banks
.macro @vic_memory_bank(bank) {
    lda vic.register.BANK_SELECT
    and #%11111100
    ora #bank
    sta vic.register.BANK_SELECT
}

// d016
.function @vic_control2_value(hor_scroll, set_38columns, multicolor) {
    .const bit3 = set_38columns ? 0 : 1
    .const bit4 = multicolor ? 1 : 0
    .return (bit4 << 4) | (bit3 << 3) | (hor_scroll & %111)
}

// d018
.function @vic_layout(font, screen) {
    .return ((font >> 10) & $0f) | ((screen >> 6) & $f0)
}

// set pointers to font mem and screen mem
.macro vic_memory_layout(font, screen) {
    lda #vic_layout(font,screen)
    sta vic.register.MEM_LAYOUT
}

.macro @vic_fill_screen(address, value) {
    lda #value
    ldx #0
!:
    sta address,x
    sta address + $100,x
    sta address + $200,x
    sta address + $300,x
    inx
    bne !-
}

.namespace sprite {
.label POINTER_OFFSET   = $07f8

.function @vic_sprite_x(nr) {
  .return $d000 + (nr << 1)
}

.function @vic_sprite_y(nr) {
  .return $d001 + (nr << 1)
}

.macro @vic_set_sprite_pointer(bank, sprite_nr, address) {
    .const BANK_OFFSET = bank * vic.memory.BANK_SIZE
    lda #(address - BANK_OFFSET) / vic.sprite.SIZE_IN_BYTES
    sta BANK_OFFSET + vic.sprite.POINTER_OFFSET + sprite_nr
}

.function @vic_sprite_pointer(bank, sprite_nr) {
   .return bank * $4000 + vic.sprite.POINTER_OFFSET + sprite_nr
}

}
