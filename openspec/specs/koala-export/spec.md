## ADDED Requirements

### Requirement: Koala Binary File Structure
The system SHALL export the image as a binary file with a total size of 10003 bytes, following the standard Koala Painter format.

#### Scenario: Correct file size
- **WHEN** an image is exported in Koala format
- **THEN** the output file size MUST be exactly 10003 bytes

### Requirement: Koala Memory Layout
The exported binary file SHALL contain the following sections in order:
1. Load Address (2 bytes): Little-endian `0x6000` (`00 60`).
2. Bitmap Data (8000 bytes): 8x8 character cell based bitmap data.
3. Screen RAM (1000 bytes): Packed colors for bit patterns `01` and `10`.
4. Color RAM (1000 bytes): Colors for bit pattern `11`.
5. Background Color (1 byte): Color for bit pattern `00`.

#### Scenario: Verify memory layout order
- **WHEN** the binary file is generated
- **THEN** the first two bytes are `00 60`
- **AND** the next 8000 bytes contain the bitmap
- **AND** the next 1000 bytes contain the screen RAM
- **AND** the next 1000 bytes contain the color RAM
- **AND** the last byte contains the background color

### Requirement: Bitmap Data Packing
The bitmap data SHALL be packed such that each byte represents 4 pixels (multi-color mode), with 2 bits per pixel.
- Bit patterns: `00`, `01`, `10`, `11`.
- Packing: MSB to LSB.
- Ordering: 8x8 character cell based (8 bytes per cell, 40 cells per row, 25 rows).

#### Scenario: Correct bitmap packing
- **WHEN** exporting a 320x200 image (effective 160x200 multi-color)
- **THEN** every 8 bytes in the bitmap section represent one 8x8 character cell (4x8 multi-color pixels)

### Requirement: Color Mapping
The system SHALL map the colors assigned to each bit pattern as follows:
- Bit pattern `00`: Uses the single background color byte.
- Bit pattern `01`: Uses the high nibble (bits 4-7) of the corresponding Screen RAM byte.
- Bit pattern `10`: Uses the low nibble (bits 0-3) of the corresponding Screen RAM byte.
- Bit pattern `11`: Uses the corresponding Color RAM byte.

#### Scenario: Verify color mapping in output
- **WHEN** a character cell has color C1 for bit pattern `01` and color C2 for bit pattern `10`
- **THEN** the corresponding Screen RAM byte MUST be `(C1 << 4) | C2`
- **AND** if bit pattern `11` has color C3, the corresponding Color RAM byte MUST be `C3`
- **AND** if bit pattern `00` has color C0, the background byte MUST be `C0`
