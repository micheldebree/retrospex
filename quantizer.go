package main

import (
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/exp/maps"
)

// PaletteDistance The distance from an RGB pixel to each key in a Palette
// In the order of the palette
type PaletteDistance map[int]float64

// ReducedPalette Palette that has been reduced to the number of bitpatterns supported in
// a specific Layer of a Retrospec
type ReducedPalette struct {
	palette     Palette
	bitpatterns map[int]int8
}

func distance(color1, color2 colorful.Color) float64 {
	return color1.DistanceRgb(color2)
	// return color1.DistanceCIE94(color2)
}

// Distances Distance from a pixel to each color in a palette
func Distances(aColor colorful.Color, palette Palette) PaletteDistance {
	result := make(PaletteDistance, len(palette))
	for palIndex, c := range palette {
		result[palIndex] = distance(aColor, c)
	}
	return result
}

// BestPixelIndex The palette index with the smallest distance
// also returns the distance itself, meaning the quantization error
// which is useful for error diffusion dithering
func BestPixelIndex(distances PaletteDistance) (index int, qerror float64) {
	bestIndex := -1
	smallestDistance := math.MaxFloat64
	for i, distance := range distances {
		if distance < smallestDistance {
			smallestDistance = distance
			bestIndex = i
		}
	}

	if bestIndex < 0 {
		panic("Could not determine best index")
	}

	return bestIndex, smallestDistance
}

func QuantizePixel(p *Pixel, pal Palette) {
	i, qerror := QuantizeToIndex(p.color, pal)
	p.paletteIndex = i
	p.quantizationError = qerror
}

func QuantizeToIndex(aColor colorful.Color, palette Palette) (int, float64) {
	return BestPixelIndex(Distances(aColor, palette))
}

func Quantize(img IndexedImage) IndexedImage {
	result := img
	for _, layer := range img.spec.layers {
		cells := getCells(result, layer)
		qCells := quantizeCells(cells, layer)
		result = combine(&qCells)
	}
	return result
}

func quantizeCells(cells []IndexedImage, layer Layer) []IndexedImage {
	result := make([]IndexedImage, len(cells))
	for ci, cell := range cells {
		result[ci] = quantizeCell(cell, layer)
	}
	return result
}

func quantizeCell(img IndexedImage, layer Layer) IndexedImage {
	// newPalette := reducePaletteKmeans(img, layer)
	newPalette := reducePalette(img, layer)

	for pi := range img.pixels {
		// pixels that are already assigned a bitpattern should not
		// be quantized as their color will not be in the reduced palette
		if !img.pixels[pi].hasBitPattern() { // has already been processed
			if layer.isLast { // last layer, all remaining pixels should be quantized against new palette
				QuantizePixel(&(img.pixels[pi]), newPalette.palette)
				img.pixels[pi].bitPattern = newPalette.bitpatterns[img.pixels[pi].paletteIndex]
			} else { // not the last layer, only process pixels that quantize to a bitpattern in the new palette
				QuantizePixel(&(img.pixels[pi]), img.palette)
				bitpattern, present := newPalette.bitpatterns[img.pixels[pi].paletteIndex]
				if present {
					img.pixels[pi].bitPattern = bitpattern
				}
			}
		}
	}
	return img
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func reducePalette(img IndexedImage, layer Layer) ReducedPalette {

	indexToCount := make(map[int]int)

	// TODO: make bitpattern map a type
	// TODO: does this make any difference?
	existingBitpatterns := make(map[int]int8)

	// count nr of pixels for each quantized color
	for _, pixel := range img.pixels {

		// pixels that are already assigned a bitpattern don't count
		if !pixel.hasBitPattern() {
			QuantizePixel(&pixel, img.palette)
			indexToCount[pixel.paletteIndex]++
		} else {
			existingBitpatterns[pixel.paletteIndex] = pixel.bitPattern
		}
	}

	// sort in reverse order of count values
	keys := maps.Keys(indexToCount)
	sort.SliceStable(keys, func(i, j int) bool {
		return indexToCount[keys[i]] > indexToCount[keys[j]]
	})

	// only keep top n
	maxColors := len(layer.bitpatterns)
	if maxColors < len(keys) {
		keys = keys[0:maxColors]
	}

	newPalette := make(Palette)
	newBitpatterns := make(map[int]int8)

	// assign bitpatterns
	i := 0
	for _, key := range keys {
		newPalette[key] = img.palette[key]
		newBitpatterns[key] = layer.bitpatterns[i]
		i++
	}

	// Add existing bitpatterns to the palette so they also get a chance
	// TODO: does this make any difference?
	for key := range existingBitpatterns {
		newPalette[key] = img.palette[key]
		newBitpatterns[key] = existingBitpatterns[key]
	}

	return ReducedPalette{newPalette, newBitpatterns}
}
