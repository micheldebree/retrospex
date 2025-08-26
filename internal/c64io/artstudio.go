package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func artstudioBinary(img *indexedimage.IndexedImage) BinaryFile {

	layer := img.Regions[0]

	nrRegions := len(layer)
	screenRam := make([]byte, nrRegions)
	for i, region := range layer {
		screenRam[i] = packScreenRamColors(region, 0b01, 0b10)
	}
	address := loadAddress(0x2000)
	bitmapChunk := getVicOrderBitmapData(img)
	screenRamChunk := BinaryChunk{"ScreenRAM", screenRam}
	paddingChunk := BinaryChunk{"Padding", make([]byte, 7)}
	return BinaryFile{address, bitmapChunk, screenRamChunk, paddingChunk}

}
