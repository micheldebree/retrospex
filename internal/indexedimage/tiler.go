package indexedimage

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"

	"golang.org/x/exp/maps"
)

func countUniqueTiles(img image.Image, spec Retrospec, pal Palette) {

	tiles := cutIntoTiles(img, spec, pal)

	theMap := make(map[string]int)
	for _, tile := range tiles {
		hash := hashTile(tile)
		count, isPresent := theMap[hash]
		if !isPresent {
			count = 0
		}
		theMap[hash] = count + 1
	}
	fmt.Printf("%v\n", theMap)
	fmt.Printf("%d unique tiles found\n", len(maps.Keys(theMap)))

}

func cutIntoTiles(img image.Image, spec Retrospec, pal Palette) []IndexedImage {

	tileW, tileH := 8, 8
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	nrRows, nrCols := h/tileH, w/tileW

	nrTiles := nrRows * nrCols
	result := make([]IndexedImage, nrTiles)
	i := 0
	for y := range nrRows {
		for x := range nrCols {

			// if img.ColorModel() != color.RGBA64Model {
			// 	panic(fmt.Sprintf("Unsupported color model: %v", img.ColorModel()))
			// }

			tile := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x*tileW, y*tileH, x*tileW+tileW, y*tileH+tileH))

			iimg := ToIndexedImage(&tile, spec, pal)

			fmt.Printf("Tile bounds: %v\n", tile.Bounds())

			qimage := Quantize(iimg)
			result[i] = qimage
			i++
		}
	}
	return result
}

func hashTile(tile IndexedImage) string {

	w, h := tile.width, tile.height

	result := make([]byte, w*h)
	i := 0
	for y := range h {
		for x := range w {
			pixel := tile.PixelAt(x, y)

			if !pixel.hasBitPattern() {
				panic("Pixel has no bit pattern!")
			}

			result[i] = byte(pixel.bitPattern)
		}
	}
	hash := sha256.New()
	hash.Write(result)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
