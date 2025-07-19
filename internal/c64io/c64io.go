package c64io

import (
	"fmt"
	"os"

	"github.com/micheldebree/retrospex/internal/indexedimage"
)

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
func reOrderToVicBitmapOrder(input []byte, bytesPerInputRow int) []byte {
	result := make([]byte, len(input))
	bytesPerOutputRow := bytesPerInputRow * 8
	numRows := len(input) / bytesPerOutputRow

	if (len(input) % bytesPerOutputRow) != 0 {
		panic("Input height should be a multiple of 8")
	}

	dstIndex := 0

	for row := range numRows {
		rowIndex := row * bytesPerOutputRow
		for col := range bytesPerInputRow {
			for byteInCol := range 8 {
				srcIndex := rowIndex + col + byteInCol*bytesPerInputRow
				result[dstIndex] = input[srcIndex]
				dstIndex++
			}
		}
	}

	return result
}

// Helper function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func SaveBinary(filename string, img *indexedimage.IndexedImage, overwrite bool) {
	bitmapData := getBitmapData(img)

	var properlyOrderedBytes []byte
	if img.Spec.BitmapByteOrder == indexedimage.VicByteOrder {
		properlyOrderedBytes = reOrderToVicBitmapOrder(bitmapData, img.Width)
	} else {
		properlyOrderedBytes = bitmapData
	}

	if !overwrite && fileExists(filename) {
		panic(fmt.Sprintf("File %s already exists", filename))
	}

	err := os.WriteFile(filename, properlyOrderedBytes, 0644)
	if err != nil {
		panic(err)
	}
}
