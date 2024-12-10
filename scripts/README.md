# Scripts

These scripts convert an image to a runnable Commodore 64 `.prg` executable.
To make non-runnable `.prg` files, leave out the `-display` argument to `png2prg`.

The scripts are dependent on a few other tools beside `retrospex`:

- [ImageMagick](https://imagemagick.org/script/index.php) for pre-processing
- [png2prg](https://github.com/staD020/png2prg) for converting to `.prg` format

## Shell scripts

Can be run with the source image as an argument. It will generate some
intermediate files that can be deleted.

## Makefiles

The `.mk` files are meant for inclusion in Gnu Make files.

They can be used like this:

- Give the source image the extension `.src.png` (for other formats, edit the
  make file), for example `mypicture.src.png`
- Include the specific `.mk` file for the format you want in the Makefile for
  your project. For example, `include source/koala.mk`
- Specify, for example, `mypicture.koala.prg` as a dependency to your own make
  targets
