package pixels

import (
	"fmt"
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// Two types of color are used:
// color.RGBA, for reading/writing Go image library images
// colorful, for internal processing

// Pixel One pixel in the image
type Pixel struct {
	X, Y  int
	Color colorful.Color
	// available after quantizing
	PaletteIndex      int
	QuantizationError float64
	// available when a bit pattern has been assigned
	BitPattern int
}

func (pixel Pixel) GetColor(palette Palette) colorful.Color {
	pixel.assertQuantized()
	return palette[pixel.PaletteIndex]
}

func (pixel Pixel) isQuantized() bool {
	return pixel.PaletteIndex >= 0
}

func (pixel Pixel) assertQuantized() {
	if !pixel.isQuantized() {
		panic(fmt.Sprintf("Pixel %v is not quantized", pixel))
	}
}

func (pixel Pixel) HasBitPattern() bool {
	return pixel.BitPattern >= 0
}

func (pixel Pixel) assertHasBitPattern() {
	if !pixel.HasBitPattern() {
		panic(fmt.Sprintf("Pixel %v does not have a bit pattern", pixel))
	}
}

func ToRGBA(aColor color.Color) color.RGBA {
	rr, gg, bb, aa := aColor.RGBA()
	return color.RGBA{R: (uint8)(rr >> 8), G: (uint8)(gg >> 8), B: (uint8)(bb >> 8), A: (uint8)(aa >> 8)}
}

func ToColor(colorful colorful.Color) color.RGBA {
	return color.RGBA{
		R: uint8(colorful.R * 256.0),
		G: uint8(colorful.G * 256.0),
		B: uint8(colorful.B * 256.0),
		A: 0xff,
	}
}

func toColorful(aColor color.RGBA) colorful.Color {
	return colorful.Color{
		R: float64(aColor.R) / 256.0,
		G: float64(aColor.G) / 256.0,
		B: float64(aColor.B) / 256.0,
	}
}

func getPixel(image *image.Image, x, y int) Pixel {
	r, g, b, a := (*image).At(x, y).RGBA()
	scale := float64(a)
	convertedColor := colorful.Color{R: float64(r) / scale, G: float64(g) / scale, B: float64(b) / scale}
	return Pixel{x, y, convertedColor, -1, -1.0, -1}
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
