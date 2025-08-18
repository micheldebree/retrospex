package pixels

import (
	"fmt"
	"image"
	"image/color"
)

// Two types of color are used:
// color.RGBA, for reading/writing Go image library images
// colorful, for internal processing

// Pixel One pixel in the image
type Pixel struct {
	X, Y  int
	Color color.Color // the color in the original image
	// available after quantizing
	PaletteIndex      int
	QuantizationError float64
	// available when a bit pattern has been assigned
	BitPattern int
}

func (pixel Pixel) GetIndexedColor(palette Palette) color.Color {
	pixel.AssertQuantized()
	return palette[pixel.PaletteIndex]
}

func (pixel Pixel) isQuantized() bool {
	return pixel.PaletteIndex >= 0
}

func (pixel Pixel) AssertQuantized() {
	if !pixel.isQuantized() {
		panic(fmt.Sprintf("Pixel %v is not quantized", pixel))
	}
}

func (pixel Pixel) HasBitPattern() bool {
	return pixel.BitPattern >= 0
}

func (pixel Pixel) AssertHasBitPattern() {
	if !pixel.HasBitPattern() {
		panic(fmt.Sprintf("Pixel %v does not have a bit pattern", pixel))
	}
}

func getPixel(image *image.Image, x, y int) Pixel {
	theColor := (*image).At(x, y)
	return Pixel{x, y, theColor, -1, -1.0, -1}
}

func GetDimensions(image *image.Image) (int, int) {
	return (*image).Bounds().Max.X - (*image).Bounds().Min.X, (*image).Bounds().Max.Y - (*image).Bounds().Min.Y
}

func GetPixels(image *image.Image) []Pixel {
	startX, startY := (*image).Bounds().Min.X, (*image).Bounds().Min.Y
	w, h := GetDimensions(image)
	result := make([]Pixel, w*h)

	for y := range h {
		offs := y * w
		for x := range w {
			result[offs+x] = getPixel(image, startX+x, startY+y)
		}
	}
	return result
}
