## ADDED Requirements

### Requirement: MCCharset Binary File Structure
The system SHALL export the multicolor charset image as a binary file. The total size is variable and depends on the number of character cells in the input image.

#### Scenario: Verify file size for full screen
- **WHEN** a 320x200 image is exported in mccharset format (squashed to 160x200)
- **THEN** the output file size MUST be exactly 9003 bytes (1000 cells * 9 bytes + 3 global colors)

#### Scenario: Verify file size for non-fullscreen
- **WHEN** an image resulting in N character cells is exported
- **THEN** the output file size MUST be (N * 9) + 3 bytes

### Requirement: MCCharset Memory Layout
The exported binary file SHALL contain the following sections in order:
1. Bitmap/Character Data (8 * N bytes): 8x8 character cell based bitmap data for N characters.
2. D800 Color RAM Data (N bytes): Individual color for each character cell (bit pattern `11`).
3. Global Colors (3 bytes): Global colors for bit patterns `00`, `01`, and `10`.

#### Scenario: Verify memory layout order
- **WHEN** the binary file is generated
- **THEN** the first 8000 bytes contain the bitmap data
- **AND** the next 1000 bytes contain the D800 color RAM
- **AND** the next 3 bytes contain the global multicolor colors

### Requirement: Bitmap Data Packing
The bitmap data SHALL be packed such that each byte represents 4 pixels (multi-color mode), with 2 bits per pixel.
- Bit patterns: `00`, `01`, `10`, `11`.
- Packing: MSB to LSB.
- Ordering: 8x8 character cell based (8 bytes per cell, 40 cells per row, 25 rows).

#### Scenario: Correct bitmap packing
- **WHEN** exporting a 320x200 image
- **THEN** every 8 bytes in the bitmap section represent one 8x8 character cell (4x8 multi-color pixels)

### Requirement: Color Mapping
The system SHALL map the colors assigned to each bit pattern as follows:
- Bit pattern `00`: Uses the first global color byte ($d021).
- Bit pattern `01`: Uses the second global color byte ($d022).
- Bit pattern `10`: Uses the third global color byte ($d023).
- Bit pattern `11`: Uses the corresponding byte from the D800 section.

#### Scenario: Verify color mapping in output
- **WHEN** an image is exported
- **THEN** the global color bytes MUST match the colors assigned to bit patterns `00`, `01`, and `10`
- **AND** the D800 section MUST contain the colors assigned to bit pattern `11` for each cell
