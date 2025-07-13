A command line application written in golang that converts any image into an
image that adheres to the limits of a Commodore 64 graphic mode.

It support different modes, like bitmap, character mode and sprites. All modes
support multicolor and hires variants.

It is focussed on applying the limitations of the Commodore 64 modes, without
doing pre-processing or actually creating a binary for use on a real Commodore 64.
It simply outputs an image that adheres to the limitations.

Pre-processing like scaling or cropping can be done with another tool and is
not in scope for Retrospex.

Creating a Commodore 64 binary format can be done with another tool like
png2prg and is not in scope. Retrospex is tuned to be used as a pre-processing
tool for png2prg, for instance by using the same palettes.

The user has a choice of different palettes. They are all the Commodore 64
standard palette, but different shades of colors depending on taste.

Popular image formats for input are supported. Output is a PNG image. There is
no limit to input image dimensions. If the user wants to make sure the image
fits on the Commodore 64 screen, they have to do the pre-processing required
themselves.
