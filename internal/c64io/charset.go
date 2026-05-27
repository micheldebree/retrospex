package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func charsetBinary(img *indexedimage.IndexedImage) BinaryFile {
	if img.Spec.Name == indexedimage.MCCharsetType {
		return mcCharsetBinary(img)
	}
	return BinaryFile{BinaryChunk{"Char data", reOrder(getBitmapData(img), img.BytesPerRow(), ByteOrders[CharsByteOrder])}}
}

func mcCharsetBinary(img *indexedimage.IndexedImage) BinaryFile {
	// 1. Bitmap Data
	bitmapData := reOrder(getBitmapData(img), img.BytesPerRow(), ByteOrders[CharsByteOrder])
	bitmapChunk := BinaryChunk{"Bitmap", bitmapData}

	// 2. D800 Data
	// Layer 1 has the char-sized regions for D800 (bit pattern 11)
	nrChars := len(img.Regions[1])
	d800Data := make([]byte, nrChars)
	for i, region := range img.Regions[1] {
		color, present := region.ColorAssignedToBitPattern(0b11)
		if present {
			d800Data[i] = byte(color)
		}
	}
	d800Chunk := BinaryChunk{"D800 RAM", d800Data}

	// 3. Global Colors (3 bytes)
	// Layer 0 has the global colors (bit patterns 00, 01, 10)
	globalColors := make([]byte, 3)
	globalRegion := img.Regions[0][0]

	c0, _ := globalRegion.ColorAssignedToBitPattern(0b00)
	c1, _ := globalRegion.ColorAssignedToBitPattern(0b01)
	c2, _ := globalRegion.ColorAssignedToBitPattern(0b10)

	globalColors[0] = byte(c0)
	globalColors[1] = byte(c1)
	globalColors[2] = byte(c2)

	globalColorsChunk := BinaryChunk{"Global colors", globalColors}

	return BinaryFile{bitmapChunk, d800Chunk, globalColorsChunk}
}
