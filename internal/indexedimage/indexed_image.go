package indexedimage

import (
	"github.com/micheldebree/retrospex/internal/pixels"
	"image"
)

// IndexedImage an image with pixels in left to right, top to bottom order
type IndexedImage struct {
	width   int
	height  int
	spec    Retrospec
	palette pixels.Palette
	pixels  []pixels.Pixel
}

// TiledImage Multiple IndexedImage tiles (all of the same size)
type TiledImage struct {
	nrRows     int
	nrCols     int
	tileWidth  int
	tileHeight int
	tiles      []IndexedImage
}

func newIndexedImage(w int, h int, spec Retrospec, pal pixels.Palette) IndexedImage {
	thePixels := make([]pixels.Pixel, w*h)
	return IndexedImage{w, h, spec, pal, thePixels}
}

func ToIndexedImage(img *image.Image, spec Retrospec, pal pixels.Palette) IndexedImage {
	thePixels := pixels.GetPixels(img)
	w, h := pixels.GetDimensions(img)
	return IndexedImage{w, h, spec, pal, thePixels}
}

func (img *IndexedImage) PixelAt(x, y int) pixels.Pixel {
	return img.pixels[y*img.width+x]
}

func (img *IndexedImage) SetPixel(pixel pixels.Pixel) {
	img.pixels[pixel.Y*img.width+pixel.X] = pixel
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
			result.SetRGBA(x, y, pixels.ToColor(pixel.GetColor(img.palette)))
		}
	}
	return result
}

// Cut up image into cells for a particular layer
func getCells(img IndexedImage, layer Layer) TiledImage {

	w, h := img.width, img.height

	nrCols, nrRows := w/layer.cellWidth, h/layer.cellHeight

	tiles := make([]IndexedImage, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {
			cell := make([]pixels.Pixel, layer.cellWidth*layer.cellHeight)
			for py := range layer.cellHeight {
				for px := range layer.cellWidth {
					cell[py*layer.cellWidth+px] = img.PixelAt(cx*layer.cellWidth+px, cy*layer.cellHeight+py)
				}
			}
			tiles[cy*nrCols+cx] = IndexedImage{layer.cellWidth, layer.cellHeight, img.spec, img.palette, cell}
		}
	}
	return TiledImage{
		nrRows,
		nrCols,
		layer.cellWidth,
		layer.cellHeight,
		tiles,
	}
}
