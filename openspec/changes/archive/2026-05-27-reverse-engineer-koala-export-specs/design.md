## Context

The Koala export functionality is already implemented in `internal/c64io/koala.go`. This project uses `openspec` to maintain a contract between requirements and implementation. Currently, the Koala export is an "untracked" feature from the perspective of `openspec`. This design document formalizes the architectural approach used by the existing implementation to ensure it serves as a stable base for the `koala-export` capability.

## Goals / Non-Goals

**Goals:**
- Formalize the data flow from `IndexedImage` to the multi-chunk `BinaryFile`.
- Document the reuse of `reOrder` for character-cell based bitmap layout.
- Ensure the bit-pattern mapping (00, 01, 10, 11) is correctly attributed to C64 memory regions.

**Non-Goals:**
- Changing the existing implementation logic.
- Adding support for other load addresses or formats in this change.

## Decisions

### Decision: Use Chunk-Based Binary Composition
The implementation uses a `BinaryFile` type, which is a slice of `BinaryChunk` objects. This allows for a clean separation of the different memory regions (Load Address, Bitmap, Screen RAM, Color RAM, Background) before they are flattened into a single byte stream.
- **Rationale**: This matches the logical structure of Commodore 64 files and makes the code easier to test and maintain.
- **Alternatives**: Manually writing to a single buffer or using a struct with fixed-size arrays. The chunk-based approach is more flexible for formats with varying section sizes.

### Decision: Leverage Centralized Reordering Logic
The `reOrder` function in `internal/c64io/binaryfile.go` is used to convert linear bitmap data into the 8x8 character cell order required by the C64 hardware.
- **Rationale**: Avoids duplicating complex indexing logic within each format-specific exporter.
- **Alternatives**: Implementing per-format reordering.

### Decision: Strict Mapping of Bit Patterns
The mapping of bit patterns to colors is hardcoded to follow the standard Koala Painter usage:
- `00` -> Background ($D021)
- `01` -> Screen RAM High Nibble
- `10` -> Screen RAM Low Nibble
- `11` -> Color RAM ($D800 range in memory, appended in file)
- **Rationale**: Ensures compatibility with all C64 software that supports the Koala format.

## Risks / Trade-offs

- **Risk: Hardcoded Load Address** → The current implementation hardcodes `0x6000`. While standard for Koala, some variants might use different addresses.
  - **Mitigation**: Keep it hardcoded for now as per the "reverse engineer existing behavior" goal. Future changes can introduce parameterization if needed.
- **Risk: Panic on Missing Background Color** → The code panics if bit pattern `00` is missing from the image's regions.
  - **Mitigation**: This is acceptable for a conversion tool where the input image is expected to conform to the target format's constraints.
