package c64io

import "fmt"

// A chunk inside a binary file
type BinaryChunk struct {
	label string // a name for reporting the output layout to the user
	data  []byte
}

type ByteOrderName int

// how the bytes are ordered in the output
type ByteOrder struct {
	ColWidthBytes, ColHeightBytes int
}

const (
	DefaultByteOrder ByteOrderName = iota
	CharsByteOrder
	SpritesByteOrder
)

var ByteOrders = map[ByteOrderName]ByteOrder{
	CharsByteOrder:   {1, 8},  // 1 byte columns, 8 byte rows
	SpritesByteOrder: {3, 21}, // 3 byte columns, 21 byte rows
}

type BinaryFile []BinaryChunk

// Get all the chunk data appended
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

// Reorders bytes to c64 bitmap byte ordering
func reOrder(input []byte, bytesPerInputRow int, byteOrder ByteOrder) []byte {

	if bytesPerInputRow%byteOrder.ColWidthBytes != 0 {
		panic(fmt.Sprintf("Bitmap data should be a multiple of %d", byteOrder.ColWidthBytes))
	}

	result := make([]byte, len(input))
	bytesPerOutputRow := bytesPerInputRow * byteOrder.ColHeightBytes
	numRows := len(input) / bytesPerOutputRow
	numCols := bytesPerInputRow / byteOrder.ColWidthBytes

	dstIndex := 0
	rowOffset := 0
	for range numRows {
		colOffset := rowOffset
		for range numCols {
			rowInColOffset := colOffset
			for range byteOrder.ColHeightBytes {
				srcIndex := rowInColOffset
				for range byteOrder.ColWidthBytes {
					result[dstIndex] = input[srcIndex]
					srcIndex++
					dstIndex++
				}
				rowInColOffset += bytesPerInputRow
			}
			colOffset += byteOrder.ColWidthBytes
		}
		rowOffset += bytesPerOutputRow
	}
	return result
}
