## Why

Add support for `mccharset` binary export format to enable usage with Commodore
64 development tools that expect a specific 9003-byte layout (bitmap data,
color RAM data, and global multicolor colors).

## What Changes

- Update `mccharset` binary export logic to generate a 9003-byte file.
- The output format will include:
    - 8000 bytes of bitmap data (character data for 1000 characters).
    - 1000 bytes of D800 data (individual color RAM for each character).
    - 3 bytes for global multicolor colors ($d021, $d022, $d023).
- Document the new binary format in the README.
- Add a unit test to verify the export format.

## Capabilities

### New Capabilities
- `mccharset-binary-export`: Implements the 9003-byte binary export format for multicolor charset mode.

### Modified Capabilities
- `koala-export`: No changes, but I reviewed it for reference.

## Impact

- `internal/c64io/charset.go`: Implementation of the new export format.
- `README.md`: Updated documentation for binary output formats.
- `internal/c64io/charset_test.go`: New test cases for verifying the export.
