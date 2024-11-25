package indexedimage

import (
	"image"
)

// IndexedImage an image with pixels in left to right, top to bottom order
type IndexedImage struct {
	width   int
	height  int
	spec    Retrospec
	palette Palette
	pixels  []Pixel
}

// TiledImage Multiple IndexedImage tiles (all of the same size)
type TiledImage struct {
	nrRows     int
	nrCols     int
	tileWidth  int
	tileHeight int
	tiles      []IndexedImage
}

func newIndexedImage(w int, h int, spec Retrospec, pal Palette) IndexedImage {
	pixels := make([]Pixel, w*h)
	return IndexedImage{w, h, spec, pal, pixels}
}

func ToIndexedImage(img *image.Image, spec Retrospec, pal Palette) IndexedImage {
	pixels := getPixels(img)
	w, h := getDimensions(img)
	return IndexedImage{w, h, spec, pal, pixels}
}

func (img *IndexedImage) PixelAt(x, y int) Pixel {
	return img.pixels[y*img.width+x]
}

func (img *IndexedImage) SetPixel(pixel Pixel) {
	img.pixels[pixel.y*img.width+pixel.x] = pixel
}

func (img *TiledImage) GetDimensions() (int, int) {
	return img.nrCols * img.tileWidth, img.nrRows * img.tileHeight
}

func combine(img *TiledImage) IndexedImage {
	pal := (*img).tiles[0].palette
	spec := (*img).tiles[0].spec
	w, h := img.GetDimensions()
	result := newIndexedImage(w, h, spec, pal)

	for _, c := range (*img).tiles {
		for _, p := range c.pixels {
			result.SetPixel(p)
		}
	}
	return result
}

// Render to a 'normal' RGBA image
func (img *IndexedImage) Render() image.Image {
	result := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: img.width, Y: img.height},
	})
	for y := range img.height {
		for x := range img.width {
			pixel := img.PixelAt(x, y)
			result.SetRGBA(x, y, toColor(pixel.getColor(img.palette)))
		}
	}
	return result
}
