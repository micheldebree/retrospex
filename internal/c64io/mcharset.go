package c64io

import "github.com/micheldebree/retrospex/internal/indexedimage"

func mcCharsetBinary(img *indexedimage.IndexedImage) BinaryFile {
	return BinaryFile{getOrderedBitmapData(img)}
}
