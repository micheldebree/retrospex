## 1. Core Logic Changes

- [x] 1.1 Update `MCCharset` template in `internal/indexedimage/retrospec_factory.go` to use 4x8 cells (instead of 8x8) to correctly represent 40 columns in 160x200 squashed input.
- [x] 1.2 Implement `mcCharsetBinary` function in `internal/c64io/charset.go` to generate the binary layout (Bitmap, D800, 3 bytes global colors).
    - Ensure it handles variable image sizes by using the actual number of regions in Layer 1.
    - Export order: Bitmap data (all chars), then D800 data (all chars), then 3 global color bytes.
- [x] 1.3 Update `charsetBinary` to call `mcCharsetBinary` specifically for `MCCharsetType`.

## 2. Documentation

- [x] 2.1 Update `README.md` to document the `mccharset` binary format.
    - Include details on how the size is calculated for non-fullscreen images: `(8 * nrChars) + (1 * nrChars) + 3`.
    - Document that oversized images are exported in full.

## 3. Testing

- [x] 3.1 Create or update `internal/c64io/charset_test.go` with a new test case for `mccharset` export.
- [x] 3.2 Verify the export produces the correct number of bytes for both fullscreen (9003 bytes) and non-fullscreen images.
- [x] 3.3 Verify correct mapping of global/local colors.
