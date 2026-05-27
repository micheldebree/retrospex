## Why

The `koala` export function is already implemented in the codebase, but its formal specification is missing from the `openspec` documentation. Formalizing this specification ensures that the existing behavior is documented as a contract, providing a clear reference for future maintenance and for implementing other Commodore 64 graphics formats.

## What Changes

- Document the memory layout of the Koala binary format (0x6000 load address, bitmap, screen RAM, color RAM, background).
- Define the mapping between image bit patterns (00, 01, 10, 11) and C64 memory locations.
- Establish the `koala-export` capability as a formal specification.

## Capabilities

### New Capabilities
- `koala-export`: Specification for exporting indexed images to the Commodore 64 Koala Painter binary format (.koa).

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec file.
     Use existing spec names from openspec/specs/. Leave empty if no requirement changes. -->

## Impact

- `internal/c64io/koala.go`: Reference implementation for the specification.
- `openspec/specs/koala-export/spec.md`: New specification file defining the format.
