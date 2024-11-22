package main

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
	x, y  int
	color colorful.Color
	// available after quantizing
	paletteIndex      int
	quantizationError float64
	// available when a bit pattern has been assigned
	bitPattern int
}

func (pixel Pixel) getColor(palette Palette) colorful.Color {
	pixel.assertQuantized()
	return palette[pixel.paletteIndex]
}

func (pixel Pixel) isQuantized() bool {
	return pixel.paletteIndex >= 0
}

func (pixel Pixel) assertQuantized() {
	if !pixel.isQuantized() {
		panic(fmt.Sprintf("Pixel %v is not quantized", pixel))
	}
}

func (pixel Pixel) hasBitPattern() bool {
	return pixel.bitPattern >= 0
}

func (pixel Pixel) assertHasBitPattern() {
	if !pixel.hasBitPattern() {
		panic(fmt.Sprintf("Pixel %v does not have a bit pattern", pixel))
	}
}

func ToRGBA(aColor color.Color) color.RGBA {
	rr, gg, bb, aa := aColor.RGBA()
	return color.RGBA{R: (uint8)(rr >> 8), G: (uint8)(gg >> 8), B: (uint8)(bb >> 8), A: (uint8)(aa >> 8)}
}

func toColor(colorful colorful.Color) color.RGBA {
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

func getDimensions(image *image.Image) (int, int) {
	return (*image).Bounds().Max.X - (*image).Bounds().Min.X, (*image).Bounds().Max.Y - (*image).Bounds().Min.Y
}

func getPixels(image *image.Image) []Pixel {
	startX, startY := (*image).Bounds().Min.X, (*image).Bounds().Min.Y
	w, h := getDimensions(image)
	result := make([]Pixel, w*h)

	for y := range h {
		offs := y * w
		for x := range w {
			result[offs+x] = getPixel(image, startX+x, startY+y)
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
			cell := make([]Pixel, layer.cellWidth*layer.cellHeight)
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
