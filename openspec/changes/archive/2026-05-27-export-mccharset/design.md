## Context

Currently, the `mccharset` mode in `retrospex` only exports character bitmap data. Commodore 64 development often requires a full binary export including color RAM (D800) and global multicolor colors. This design addresses the need for a 9003-byte export format.

## Goals / Non-Goals

**Goals:**
- Implement a multi-chunk binary export for `mccharset` mode (Bitmap + D800 + 3 global colors).
- Correct the `MCCharset` template to use proper 4x8 character cells (for 160x200 squashed input).
- Handle variable image sizes by exporting only the characters present in the input (no padding).
- Update `README.md` with the new format details and size calculation.
- Add unit tests.

**Non-Goals:**
- Supporting other charset variations (e.g. 2048-byte charsets) in this specific change.
- Modifying `sccharset` or other modes.

## Decisions

### 1. Update `MCCharset` Template
The current `MCCharset` template uses 8x8 cells for a mode that squashes input to 160 pixels wide. This results in only 20 columns. It will be changed to 4x8 cells to correctly represent 40 columns on the C64.

### 2. Multi-chunk Binary Export
`charsetBinary` in `internal/c64io/charset.go` will be updated to handle `MCCharsetType` specifically, returning a `BinaryFile` with three chunks:
1. **Bitmap Data** (8000 bytes): Character data for 1000 characters, reordered to C64 character cell format.
2. **D800 Data** (1000 bytes): Extracted from the second layer's regions (bit pattern `11`).
3. **Global Colors** (3 bytes): Extracted from the first layer's region (bit patterns `00`, `01`, `10`).

### 3. Extracting Colors from Regions
- Global colors will be fetched from `img.Regions[0][0]` using `ColorAssignedToBitPattern`.
- D800 colors will be fetched from each region in `img.Regions[1]` for bit pattern `0b11`.

## Risks / Trade-offs

- **[Risk]** Existing users of `mccharset` mode might expect only the bitmap data.
- **[Mitigation]** The change is specifically requested for "mccharset" export. If users need only bitmap data, they can truncate the file or we can consider a separate flag in the future, but for now, we follow the "full" format.
