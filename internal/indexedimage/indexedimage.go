package indexedimage

import (
	"fmt"
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

// IndexedImage an image with pixels in left to right, top to bottom order
type IndexedImage struct {
	Width   int
	Height  int
	Spec    Retrospec
	Palette pixels.Palette
	Pixels  []pixels.Pixel
	Regions map[int][]Region // regions per layer
}

func ToIndexedImage(img *image.Image, spec Retrospec, pal pixels.Palette) IndexedImage {
	thePixels := pixels.GetPixels(img)
	w, h := pixels.GetDimensions(img)

	regions := make(map[int][]Region, len(spec.Layers))

	result := IndexedImage{w, h, spec, pal, thePixels, regions}

	// for each layer, create the regions
	for layerIndex, layer := range spec.Layers {
		regions[layerIndex] = getRegions(&result, w, h, layer)
		fmt.Printf("Created %d regions for layer %d with bit patterns %v (cell size: %dx%d)\n",
			len(regions[layerIndex]), layerIndex, layer.Bitpatterns, layer.CellWidth, layer.CellHeight)
	}

	return result
}

// Cut up image into regions for a particular layer
func getRegions(img *IndexedImage, w, h int, layer Layer) []Region {

	nrCols, nrRows := w/layer.CellWidth, h/layer.CellHeight

	regions := make([]Region, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {

			regionIndex := cy*nrCols + cx

			regions[regionIndex] = Region{
				img,
				cx * layer.CellWidth,
				cy * layer.CellHeight,
				layer.CellWidth,
				layer.CellHeight,
				make(map[int]int),
				make(map[int]int),
				layer.IsLast,
			}
			regions[regionIndex].initBitPatterns(layer)
		}
	}
	return regions
}

func (img *IndexedImage) PixelAt(x, y int) pixels.Pixel {
	return img.Pixels[y*img.Width+x]
}

func (img *IndexedImage) SetPixel(pixel pixels.Pixel) {
	img.Pixels[pixel.Y*img.Width+pixel.X] = pixel
}

// Render to a 'normal' RGBA image
func (img *IndexedImage) Render() image.Image {

	pixelWidth := img.Spec.BitsPerPixel
	resultWidth, resultHeight := img.Width*pixelWidth, img.Height

	result := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: resultWidth, Y: resultHeight},
	})
	for y := range img.Height {
		for x := range img.Width {
			pixel := img.PixelAt(x, y)
			resultX := x * pixelWidth
			for p := range pixelWidth {
				result.Set(resultX+p, y, pixel.GetIndexedColor(img.Palette))
			}
		}
	}
	return result
}

// BytesPerRow TODO: this is not true for sprites
func (img *IndexedImage) BytesPerRow() int {
	return img.Width / img.Spec.PixelsPerByte()
}
