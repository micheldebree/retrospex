# Issues and features

## Next release

- [ ] Store c64 binary formats (bitmap and color data)
    - [x] Koala
    - [ ] Hires
    - [ ] Mixed Hires
    - [ ] Characters
    - [ ] Mixed Characters
    - [ ] Sprites
- [x] Override background color (or more generic, force bitpattern colors)
- [ ] Better error handling
- [ ] Speedup
- [x] Different color spaces

## Feature requests

- [ ] Do not use panic for validations
- [ ] Re-use characters
- [ ] Error diffusion dithering
- [ ] Process two pixels from the original as one, so we don't have to halve
  multicolor in preprocessing
- [ ] K-means clustering for reducing palette
- [ ] Different algorithms for measuring distance
- [ ] Parallel processing of images in folder (for animation)

## References

- <https://dev.to/bytebodger/determining-the-rgb-distance-between-two-colors-4n91>
