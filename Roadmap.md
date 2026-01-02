# Issues and features

## Branches

- feature/save-binary
  - [ ] Fix sprite modes
  - [ ] Test Artstudio mode
  - [ ] Test Koala mode
  - [ ] Test MC chars mode
  - [ ] Test Hires chars mode
- Java
  - Java conversion
  - Doesn't work yet
- Rust
  - Doesn't work yet

## Next release

- [ ] Store c64 binary formats (bitmap and color data)
  - [x] Koala
  - [x] Hires (Artstudio)
  - [x] Mixed Hires
  - [ ] Characters
  - [ ] Mixed Characters
  - [ ] Sprites
- [ ] Simple resizer
  - [ ] Remove scripts
- [ ] Handle images with non-byte widths
- [x] Override background color (or more generic, force bitpattern colors)
- [x] Different color spaces
- [ ] Document binary formats

## Feature requests

- [ ] Multipaint export
- [ ] ChatGPT recommendations
- [ ] PALette
- [ ] Speedup when using plain RGB
- [ ] Do not use panic for validations
- [ ] Re-use characters --> separate tool?
- [ ] Error diffusion dithering
- [-] Process two pixels from the original as one, so we don't have to halve
  multicolor in preprocessing
- [-] (K-means) clustering for reducing palette
- [ ] Parallel processing of images in folder (for animation)
- [ ] Custom byte orders?

## References

- <https://dev.to/bytebodger/determining-the-rgb-distance-between-two-colors-4n91>
