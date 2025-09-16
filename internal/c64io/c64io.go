package c64io

import (
	"fmt"
	"os"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/io"
)

type BinaryChunk struct {
	label string
	data  []byte
}

type ByteOrderName string

type ByteOrder struct {
	ColWidthBytes, ColHeightBytes int
}

const (
	DefaultByteOrder ByteOrderName = "default"
	CharsByteOrder   ByteOrderName = "chars"
	SpritesByteOrder ByteOrderName = "sprites"
)

var ByteOrders = map[ByteOrderName]ByteOrder{
	DefaultByteOrder: ByteOrder{1, 1},
	CharsByteOrder:   ByteOrder{1, 8},
	SpritesByteOrder: ByteOrder{3, 21},
}

type BinaryFile []BinaryChunk

var binaryFactories = map[indexedimage.RetrospecName]func(*indexedimage.IndexedImage) BinaryFile{
	indexedimage.MCCharsetType: charsetBinary,
	indexedimage.SCCharsetType: charsetBinary,
	indexedimage.KoalaType:     koalaBinary,
	indexedimage.HiresType:     artstudioBinary,
}

func (binaryFile BinaryFile) getBytes() []byte {
	result := make([]byte, 0)
	for _, chunk := range binaryFile {
		result = append(result, chunk.data...)
	}
	return result
}

func (binaryFile BinaryFile) printLayout() {
	for _, chunk := range binaryFile {
		fmt.Printf("%s: %d bytes\n", chunk.label, len(chunk.data))
	}
}

// Bitpatterns are packed into bytes;
// if the image's spec has a bitPatternSize of 1, 8 bits are packed into one byte, msb to lsb order
// if the image's spec has a bitPatternSize of 2, 4 bit patterns are packed into one byte, msb to lsb order
func getBitmapData(img *indexedimage.IndexedImage) []byte {
	width, height := img.Width, img.Height
	bitPatternSize := img.Spec.BitPatternSize
	pixelsPerByte := img.Spec.PixelsPerByte()

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
			pixel.AssertHasBitPattern()

			outIndex := byteIndex / pixelsPerByte
			shiftLeft := bitPatternSize * (pixelsPerByte - 1 - byteIndex%pixelsPerByte)

			result[outIndex] |= byte(pixel.BitPattern) << shiftLeft
			byteIndex++
		}
	}
	return result
}

// Reorders bytes to c64 bitmap byte ordering
func reOrder(input []byte, bytesPerInputRow int, byteOrder ByteOrder) []byte {

	// blocks are 1 byte wide, 8 bytes high (character layout)
	// TODO: 3 bytes wide, 21 bytes high (sprite layout)
	// colWidthBytes, colHeightBytes := 1, 8

	result := make([]byte, len(input))
	bytesPerOutputRow := bytesPerInputRow * byteOrder.ColHeightBytes
	numRows := len(input) / bytesPerOutputRow
	numCols := bytesPerInputRow / byteOrder.ColWidthBytes

	dstIndex := 0

	for row := range numRows {
		rowOffset := row * bytesPerOutputRow
		for col := range numCols {
			colOffset := rowOffset + col*byteOrder.ColWidthBytes
			for rowInCol := range byteOrder.ColHeightBytes {
				rowInColOffset := colOffset + rowInCol*bytesPerInputRow
				for colInCol := range byteOrder.ColWidthBytes {
					srcIndex := rowInColOffset + colInCol
					result[dstIndex] = input[srcIndex]
					dstIndex++
				}
			}
		}
	}
	return result
}

// Make a chunk of two bytes denoting the load address
func loadAddress(address uint16) BinaryChunk {
	return BinaryChunk{"Load address", []byte{byte(address & 0xff), byte(address >> 8)}}
}

// pack two colors associated with a bitpattern into one screenRam byte
func packScreenRamColors(region indexedimage.Region, upperBitpattern, lowerBitpattern int) byte {
	var result byte
	upperNibble, present := region.ColorAssignedToBitPattern(upperBitpattern)
	if present {
		result = byte(upperNibble) << 4
	}

	lowerNibble, present := region.ColorAssignedToBitPattern(lowerBitpattern)
	if present {
		result |= byte(lowerNibble & 0x0f)
	}
	return result
}

func SaveBinary(filename string, img *indexedimage.IndexedImage, overwrite bool) {

	binaryFactory, present := binaryFactories[img.Spec.Name]

	if !present {
		panic(fmt.Sprintf("Binary export for %s is not supported.", img.Spec.Name))
	}

	data := binaryFactory(img)

	io.AssertOverwrite(filename, overwrite)

	err := os.WriteFile(filename, data.getBytes(), 0644)
	if err != nil {
		panic(err)
	}

	data.printLayout()
}
