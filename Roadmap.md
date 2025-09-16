# Issues and features

## Next release

- [ ] Store c64 binary formats (bitmap and color data)
  - [x] Koala
  - [x] Hires
  - [ ] Mixed Hires
  - [ ] Characters
  - [ ] Mixed Characters
  - [ ] Sprites
- [ ] Handle images with non-byte widths
- [x] Override background color (or more generic, force bitpattern colors)
- [ ] Speedup when using plain RGB
- [x] Different color spaces
- [ ] Document binary formats
- [ ] PALette

## Feature requests

- [ ] Do not use panic for validations
- [ ] Re-use characters --> separate tool?
- [ ] Error diffusion dithering
- [ ] Process two pixels from the original as one, so we don't have to halve
      multicolor in preprocessing
- [ ] (K-means) clustering for reducing palette
- [ ] Parallel processing of images in folder (for animation)
- [ ] Custom byte orders?

## References

- <https://dev.to/bytebodger/determining-the-rgb-distance-between-two-colors-4n91>
