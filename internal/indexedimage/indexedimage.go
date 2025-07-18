package indexedimage

import (
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

// IndexedImage an image with pixels in left to right, top to bottom order
type IndexedImage struct {
	Width   int
	Height  int
	Spec    Retrospec
	palette pixels.Palette
	pixels  []pixels.Pixel
}

func ToIndexedImage(img *image.Image, spec Retrospec, pal pixels.Palette) IndexedImage {
	thePixels := pixels.GetPixels(img)
	w, h := pixels.GetDimensions(img)
	return IndexedImage{w, h, spec, pal, thePixels}
}

func (img *IndexedImage) PixelAt(x, y int) pixels.Pixel {
	return img.pixels[y*img.Width+x]
}

func (img *IndexedImage) SetPixel(pixel pixels.Pixel) {
	img.pixels[pixel.Y*img.Width+pixel.X] = pixel
}

// Render to a 'normal' RGBA image
func (img *IndexedImage) Render() image.Image {
	result := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: img.Width, Y: img.Height},
	})
	for y := range img.Height {
		for x := range img.Width {
			pixel := img.PixelAt(x, y)
			result.SetRGBA(x, y, pixels.ToColor(pixel.GetColor(img.palette)))
		}
	}
	return result
}
