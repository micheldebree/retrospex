# retrospex

Converts an image to match Commodore 64 graphic modes specifications.
Designed to be used as part of a development toolchain, for example together with other tools like:

- [ImageMagick](https://imagemagick.org/index.php)
- [png2prg](https://github.com/staD020/png2prg). retrospex is specifically designed for compatibility with `png2prg`
- [spot](https://github.com/spartaomg/spot_cpp)

## Example toolchain

```mermaid
graph TD
image((image)) --> ImageMagick
ImageMagick -- color adjusted, cropped --> retrospex
retrospex -- quantized --> png2prg
retrospex -- quantized koala --> spot
c64image((C64 binary))
png2prg --> c64image
spot --> c64image
```

## Usage

```bash
Usage: retrospex [options] input.png

Options:

	-o
		Output filename (default out.png)
	-m
		Graphics mode. (default koala), One of koala,hires,mixedcharset,mccharset,sccharset
	-p
		Palette (default colodore). One of archmage,electric,colodore,vice,vice_old_lum,ste,perplex_1,perplex_2,lenyn,pepto,vice_371_internal,pantaloon
	-dm
		Dither matrix (default bayer4x4). One of none,bayer2x2,bayer4x4,bayer8x8
	-dd
		Dither depth (default 25). 0-255
```
