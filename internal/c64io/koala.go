package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func koalaBinary(img *indexedimage.IndexedImage) BinaryFile {

	// background is  layer 0, one region, bit pattern 00
	backgroundColor := img.Regions[0][0].BitpatternToColor[0b00]

	// layer 1 is char-sized regions, bit patterns 01 10 and 11
	nrRegions := len(img.Regions[1])
	screenRam := make([]byte, nrRegions)
	colorRam := make([]byte, nrRegions)
	for i, region := range img.Regions[1] {

		screenRam[i] = packScreenRamColors(region, 0b01, 0b10)

		colorByte, present := region.BitpatternToColor[0b11]
		if present {
			colorRam[i] = byte(colorByte)
		}

	}

	address := loadAddress(0x6000)
	bitmapChunk := getVicOrderBitmapData(img)
	screenRamChunk := BinaryChunk{"ScreenRAM", screenRam}
	colorRamChunk := BinaryChunk{"ColorRAM", colorRam}
	backgroundChunk := BinaryChunk{"Background color", []byte{byte(backgroundColor)}}

	return BinaryFile{address, bitmapChunk, screenRamChunk, colorRamChunk, backgroundChunk}
}
