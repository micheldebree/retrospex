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
    .label COLOR_RAM        = $d800
    .label BANK_SELECT      = $dd00
}

// some common VIC data sizes
size: {
    .label SPRITE           = $40
    .label FONT             = $800
    .label BITMAP           = $2000
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

 // RSEL|  Display window height   | First line  | Last line
 // ----+--------------------------+-------------+----------
 //   0 | 24 text lines/192 pixels |   55 ($37)  | 246 ($f6)
 //   1 | 25 text lines/200 pixels |   51 ($33)  | 250 ($fa)
 //
 // CSEL|   Display window width   | First X coo. | Last X coo.
 // ----+--------------------------+--------------+------------
 //   0 | 38 characters/304 pixels |   31 ($1f)   |  334 ($14e)
 //   1 | 40 characters/320 pixels |   24 ($18)   |  343 ($157)

// d011

// 3.7.3.1. Standard text mode (ECM/BMM/MCM=0/0/0)
// 3.7.3.2. Multicolor text mode (ECM/BMM/MCM=0/0/1)
// 3.7.3.3. Standard bitmap mode (ECM/BMM/MCM=0/1/0)
// 3.7.3.4. Multicolor bitmap mode (ECM/BMM/MCM=0/1/1)
// 3.7.3.5. ECM text mode (ECM/BMM/MCM=1/0/0)
// 3.7.3.6. Invalid text mode (ECM/BMM/MCM=1/0/1)
// 3.7.3.7. Invalid bitmap mode 1 (ECM/BMM/MCM=1/1/0)
// 3.7.3.8. Invalid bitmap mode 2 (ECM/BMM/MCM=1/1/1)


// | $d011 |RST8| ECM| BMM| DEN|RSEL|    YSCROLL (3)  
.function @vic_d011(ver_scroll, set_24lines, ecm, bitmap) {

    .const bit6 = ecm ? 1 : 0
    .const bit5 = bitmap ? 1 : 0
    .const bit3 = set_24lines ? 0 : 1
    .return %00010000 | bit6 << 6 | bit5 << 5 | bit3 << 3 | ver_scroll & %00000111
}

.function @vic_enable(value, bits) {
    .return value | bits
}

.function @vic_disable(value, bits) {
    .return value & bits^%11111111
}

// d016
.function @vic_control2_value(hor_scroll, set_38columns, multicolor) {
    .const bit3 = set_38columns ? 0 : 1
    .const bit4 = multicolor ? 1 : 0
    .return (bit4 << 4) | (bit3 << 3) | (hor_scroll & %111)
}

// d018

//  $d018 |VM13|VM12|VM11|VM10|CB13|CB12|CB11| - |    Memory pointers
// VM = video matrix
// CB = characters / bitmap (bitmap only CB13)


.function @vic_layout(bitmap, screen) {
    .return ((bitmap >> 10) & $0f) | ((screen >> 6) & $f0)
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
