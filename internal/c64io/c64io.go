package c64io

import (
	"fmt"
	"os"

	"github.com/micheldebree/retrospex/internal/indexedimage"
)

type BinaryChunk struct {
	label string
	data  []byte
}

type BinaryFile []BinaryChunk

var binaryFactories = map[indexedimage.RetrospecName]func(*indexedimage.IndexedImage) BinaryFile{
	indexedimage.MCCharsetType: mcCharsetBinary,
	indexedimage.KoalaType:     koalaBinary,
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
// if the image's spec has a bitPatternSize of 1, 8 bit patterns are packed into one byte, msb to lsb order
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
// The input is row-first, with bytesPerRow bytes per row
// The output is row-first for a group of 8 bytes, but column first within each group of 8 bytes
func reOrderToVicBitmapOrder(input []byte, bytesPerInputRow int) []byte {
	result := make([]byte, len(input))
	bytesPerOutputRow := bytesPerInputRow * 8
	numRows := len(input) / bytesPerOutputRow

	dstIndex := 0

	for row := range numRows {
		rowIndex := row * bytesPerOutputRow
		for col := range bytesPerInputRow {
			colIndex := rowIndex + col
			for byteInCol := range 8 {
				srcIndex := colIndex + byteInCol*bytesPerInputRow
				result[dstIndex] = input[srcIndex]
				dstIndex++
			}
		}
	}
	return result
}

// get the bitmap data, in the right ordering
func getOrderedBitmapData(img *indexedimage.IndexedImage) BinaryChunk {

	bitmapData := getBitmapData(img)

	if img.Spec.BitmapByteOrder == indexedimage.VicByteOrder {
		return BinaryChunk{"Bitmap (vic order)", reOrderToVicBitmapOrder(bitmapData, img.BytesPerRow())}
	}
	return BinaryChunk{"Bitmap (normal order)", bitmapData}
}

// Helper function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func mcCharsetBinary(img *indexedimage.IndexedImage) BinaryFile {
	return BinaryFile{getOrderedBitmapData(img)}
}

func koalaBinary(img *indexedimage.IndexedImage) BinaryFile {

	// background is  layer 0, one region, bit pattern 00
	backgroundColor := img.Regions[0][0].BitpatternToColor[0b00]

	// layer 1 is char-sized regions, bit patterns 01 10 and 11
	nrRegions := len(img.Regions[1])
	screenRam := make([]byte, nrRegions)
	colorRam := make([]byte, nrRegions)
	for i, region := range img.Regions[1] {

		var screenByte byte
		upperNibble, present := region.BitpatternToColor[0b10]
		if present {
			screenByte = byte(upperNibble) << 4
		}

		lowerNibble, present := region.BitpatternToColor[0b10]
		if present {
			screenByte |= byte(lowerNibble)
		}

		colorByte, present := region.BitpatternToColor[0b11]
		if present {
			colorRam[i] = byte(colorByte)
		}

		screenRam[i] = screenByte

	}

	bitmapChunk := getOrderedBitmapData(img)
	screenRamChunk := BinaryChunk{"ScreenRAM", screenRam}
	colorRamChunk := BinaryChunk{"ColorRAM", colorRam}
	backgroundChunk := BinaryChunk{"Background color", []byte{byte(backgroundColor)}}

	return BinaryFile{bitmapChunk, screenRamChunk, colorRamChunk, backgroundChunk}

}

func SaveBinary(filename string, img *indexedimage.IndexedImage, overwrite bool) {

	binaryFactory, present := binaryFactories[img.Spec.Name]

	if !present {
		panic(fmt.Sprintf("Binary export for %s is not supported.", img.Spec.Name))
	}

	data := binaryFactory(img)

	saveData(filename, data, overwrite)
	data.printLayout()
}

func saveData(filename string, data BinaryFile, overwrite bool) {

	if !overwrite && fileExists(filename) {
		panic(fmt.Sprintf("File %s already exists", filename))
	}

	err := os.WriteFile(filename, data.getBytes(), 0644)
	if err != nil {
		panic(err)
	}
}
