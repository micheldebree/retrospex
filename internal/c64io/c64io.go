package c64io

import (
	"fmt"
	"github.com/micheldebree/retrospex/internal/indexedimage"
)

// This function returns an array of bytes representing the bitpatterns assigned to the pixels in the image.
// If a pixel has no bitpattern assigned, it panics
// Bitpatterns are packed into bytes;
// if the image's spec has a bitPatternSize of 1, 8 bit patterns are packed into one byte, msb to lsb order
// if the image's spec has a bitPatternSize of 2, 4 bit patterns are packed into one byte, msb to lsb order
func getBitmapData(img *indexedimage.IndexedImage) []byte {
	width, height := img.Width, img.Height
	bitPatternSize := img.Spec.BitPatternSize
	pixelsPerByte := 8 / bitPatternSize

	if width%pixelsPerByte != 0 {
		panic(fmt.Sprintf("Image width should be a multiple of %d", pixelsPerByte))
	}

	// Calculate the number of bytes needed
	numPixels := width * height
	numBytes := (numPixels + pixelsPerByte - 1) / pixelsPerByte // Ceiling division

	result := make([]byte, numBytes)
	byteIndex := 0

	for y := range height {
		for x := range width {
			pixel := img.PixelAt(x, y)
			pixel.AssertHasBitPattern() // Ensure the pixel has a bit pattern

			outIndex := byteIndex / pixelsPerByte

			shiftLeft := bitPatternSize * (pixelsPerByte - 1 - byteIndex%pixelsPerByte)

			result[outIndex] |= byte(pixel.BitPattern) << shiftLeft
			byteIndex++
		}
	}

	return result
}

// Reorders bytes to c64 bitmap byte ordering
// The input is row-first, with bytesPerRow bytes per row
// The output is row-first for a group of 8 bytes, but column first within each group of 8 bytes
func reOrderToC64BitmapOrder(input []byte, bytesPerInputRow int) []byte {

	result := make([]byte, len(input))

	bytesPerOutputRow := bytesPerInputRow * 8 // bytes per 8-byte high row in the output

	numRows := len(input) / bytesPerOutputRow // number of 8-byte high rows

	if (len(input) % bytesPerOutputRow) != 0 {
		panic("Input height should be a multiple of 8")
	}

	dstIndex := 0

	// walk over the input array in c64 order
	for row := range numRows { // each 8-byte high row
		rowIndex := row * bytesPerOutputRow // The start index (in the input) of this row
		for col := range bytesPerInputRow {
			colIndex := rowIndex + col // start index of this column
			for byteInCol := range 8 {
				srcIndex := colIndex + byteInCol*bytesPerInputRow
				result[dstIndex] = input[srcIndex]
				dstIndex++
			}
		}
	}

	return result
}
