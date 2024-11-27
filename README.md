# retrospex

Converts an image to match Commodore 64 graphic modes specifications.

Retrospex only does two things:

- Quantize the image colors to the 16-color Commodore 64 palette
- Apply restrictions that various graphic modes impose on the number of unique colors
  used in a specific area

Things like cropping/scaling to the right size, or making sure multicolor
pixels are two pixels wide is not part of this tool and can be done with other
tools in a pre-processing step.

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

	-o
		Output filename (default out.png)
	-m
		Graphics mode. (default koala), One of hires,mixedcharset,mccharset,scccharset,mcibitmap,scsprites,mcsprites,koala
	-p
		Palette (default colodore). One of archmage,electric,colodore,pepto,vice,vice_old_lum,vice_371_internal,pantaloon,ste,perplex_1,lenyn,perplex_2
	-dm
		Dither matrix (default bayer4x4). One of bayer4x4,bayer8x8,none,bayer2x2
	-dd
		Dither depth (default 25). 0-255
```

## Scripts

The `/scripts` folder contains examples on how to script a build from
source image to Commodore 64 file. Example makefiles are included aswell.
