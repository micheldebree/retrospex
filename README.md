# retrospex

Converts an image to match Commodore 64 graphic modes specifications.

Retrospex only does a few things:

- Quantize the image colors to the 16-color Commodore 64 palette
- Apply restrictions that various graphic modes impose on the number of unique colors
  used in a specific area
- Save a PNG image for further processing, or save binary data for use in Commodore 64 development

Designed to be used as part of a development toolchain, for example:

- [ImageMagick](https://imagemagick.org/index.php) (process an image, like
  resizing and cropping)
- [png2prg](https://github.com/staD020/png2prg). (convert to a Commodore 64
  compatible format) retrospex is specifically designed for compatibility with `png2prg`
- [spot](https://github.com/spartaomg/spot_cpp) (convert and optimize to
  Commodore 64 compatible format)

## Usage

```bash

Usage: retrospex [options] input.png

Options:

  -bpc string
     Force bitpattern/color pairs.
     For example 0:0 to force background black.
     For example 0:0,1:1,2:15,3:13 to force colors for all 4 bitpatterns
  -cs string
     Colorspace conversion.
     One of rgb (no conversion),lineairRgb,xyz,xyy,lab (default "lineairRgb")
  -dd int
     Dither depth (0-255). Depth of dithering. (default 25)
  -dm string
     A preset ordered dithering matrix.
     One of none,bayer2x2,bayer4x4,bayer8x8 (default "bayer4x4")
  -f string
     Output format.
     png for an image to convert for instance with png2prg
     bin for binary data to use in Commodore64 development (see documentation for structure) (default "png")
  -m string
     Graphics mode.
     One of koala,hires,mixedhires,mixedcharset,mccharset,sccharset,mcibitmap,scsprites,mcsprites (default "koala")
  -o string
     Output filename (default "out.png")
  -overwrite
     Allow overwriting output file
  -p string
     Palette.
     One of colodore,pepto,vice,vice_old_lum,vice_371_internal,pantaloon,archmage,electric,ste,perplex_1,perplex_2,lenyn (default "colodore")
  -resize
     Resize input image to fill target dimensions

```

## Notes

- Make sure input images are properly sized and pre-processed if necesary.
  `-resize` is merely included as a convenience, and only supports resizing to
  full screen. So it probably doesn't make much sense for charset and sprite
  modes.

## Binary output

There are no restrictions on size, which means that binary data size may vary.
Below are the sizes for a typical binary file for each mode.

### Koala (multicolor bitmap)

- 2 bytes loading address
- 8000 bytes bitmap data
- 1000 bytes of screen RAM
- 1000 bytes of color RAM
- 1 byte of background color

### Hires bitmap

- 2 bytes loading address
- 8000 bytes bitmap data
- 1000 bytes of screen RAM
- 7 bytes of padding (to make it compatible with Art Studio format)

### sccharset (single color charset)

- n x 8 bytes of character data for n characters.

### mccharset (multicolor charset)

The output format consists of three parts:
- **Bitmap Data** (8 * n bytes): 8x8 character cell based bitmap data for n characters.
- **D800 Data** (n bytes): Individual color RAM for each character (extracted from bit pattern `11`).
- **Global Colors** (3 bytes): Global multicolor colors (extracted from bit patterns `00`, `01`, and `10`).

The total size is `(9 * n) + 3` bytes. For a fullscreen 1000-character image (40x25 characters), the size is exactly 9003 bytes.
Oversized images are not truncated; all characters present in the input image are exported.
