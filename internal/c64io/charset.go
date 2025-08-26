package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func charsetBinary(img *indexedimage.IndexedImage) BinaryFile {
	return BinaryFile{getVicOrderBitmapData(img)}
}
