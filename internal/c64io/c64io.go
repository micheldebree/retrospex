package c64io

import (
	"github.com/micheldebree/retrospex/internal/indexedimage"
)

// This function returns an array of bytes representing the bitpatterns assigned to the pixels in the image.
// If a pixel has no bitpattern assigned, it panics
// Bitpatterns are packed into bytes;
// if the image's spec has a bitPatternSize of 1, 8 bit patterns are packed into one byte, msb to lsb order
// if the image's spec has a bitPatternSize of 2, 4 bit patterns are packed into one byte, msb to lsb order
// The order of the bytes in the files is the specific Commodore 64 layout of bytes in bitmap mode
func getBitmapData(img *indexedimage.IndexedImage) {
}
