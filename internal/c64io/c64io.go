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
func getBitmapData(img *indexedimage.IndexedImage) []byte {
    width, height := img.width, img.height
    bitPatternSize := img.spec.bitsPatternSize

    // Calculate the number of bytes needed
    numBitPatterns := width * height
    numBytes := (numBitPatterns + 7) / 8 // Ceiling division

    result := make([]byte, numBytes)
    byteIndex := 0
    bitPosition := 0

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            pixelIndex := y*width + x
            pixel := &img.pixels[pixelIndex]
            pixel.assertHasBitPattern() // Ensure the pixel has a bit pattern

            bitPattern := pixel.BitPattern
            if bitPatternSize == 1 {
                result[byteIndex] |= byte(bitPattern) << (7 - bitPosition)
                bitPosition++
                if bitPosition == 8 {
                    bitPosition = 0
                    byteIndex++
                }
            } else if bitPatternSize == 2 {
                result[byteIndex] |= byte(bitPattern) << (6 - bitPosition)
                bitPosition += 2
                if bitPosition == 8 {
                    bitPosition = 0
                    byteIndex++
                }
            } else {
                panic("Unsupported bit pattern size")
            }
        }
    }

    return result
}
