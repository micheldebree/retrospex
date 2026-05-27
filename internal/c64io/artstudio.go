package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func artstudioBinary(img *indexedimage.IndexedImage) BinaryFile {

	// There is only one layer
	if len(img.Regions) != 1 {
		panic("Image should have 1 layer.")
	}
	layer := img.Regions[0]

	// screen cells are always 8x8
	nrScreenCells := (img.Width / 8) * (img.Height / 8)
	screenRam := make([]byte, nrScreenCells)

	for i, region := range layer {
		screenRam[i] = packScreenRamColors(region, 1, 0)
	}

	address := loadAddress(0x2000)
	bitmapChunk := BinaryChunk{"Bitmap", reOrder(getBitmapData(img), img.BytesPerRow(), ByteOrders[CharsByteOrder])}
	screenRamChunk := BinaryChunk{"ScreenRAM", screenRam}
	paddingChunk := BinaryChunk{"Padding", make([]byte, 7)}
	return BinaryFile{address, bitmapChunk, screenRamChunk, paddingChunk}

}
