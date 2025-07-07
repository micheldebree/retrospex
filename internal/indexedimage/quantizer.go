package indexedimage

import (
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/micheldebree/retrospex/internal/pixels"
	"golang.org/x/exp/maps"
)

// PaletteDistance The distance from an RGB pixel to each key in a Palette
// In the order of the palette
type PaletteDistance map[int]float64

// ReducedPalette Palette that has been reduced to the number of bitpatterns supported in
// a specific Layer of a Retrospec
type ReducedPalette struct {
	palette     pixels.Palette
	bitpatterns map[int]int
}

func distance(color1, color2 colorful.Color) float64 {
	// TODO: make an option
	return color1.DistanceLinearRGB(color2)
}

// distances Distance from a pixel to each color in a palette
func distances(aColor colorful.Color, palette pixels.Palette) PaletteDistance {
	result := make(PaletteDistance, len(palette))
	for palIndex, c := range palette {
		result[palIndex] = distance(aColor, c)
	}
	return result
}

// bestPixelIndex The palette index with the smallest distance
// also returns the distance itself, meaning the quantization error
// which is useful for error diffusion dithering
func bestPixelIndex(distances PaletteDistance) (index int, qerror float64) {
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

func quantizePixel(p *pixels.Pixel, pal pixels.Palette) {
	i, qerror := QuantizeToIndex(p.Color, pal)
	p.PaletteIndex = i
	p.QuantizationError = qerror
}

func QuantizeToIndex(aColor colorful.Color, palette pixels.Palette) (int, float64) {
	return bestPixelIndex(distances(aColor, palette))
}

// Tile represents a region of an IndexedImage
type Tile struct {
	img           *IndexedImage
	x, y          int
	width, height int
}

// Cut up image into tiles for a particular layer
func getTiles(img IndexedImage, layer Layer) []Tile {
	w, h := img.width, img.height

	nrCols, nrRows := w/layer.cellWidth, h/layer.cellHeight

	tiles := make([]Tile, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {
			tiles[cy*nrCols+cx] = Tile{
				&img,
				cx * layer.cellWidth,
				cy * layer.cellHeight,
				layer.cellWidth,
				layer.cellHeight,
			}
		}
	}
	return tiles
}

func quantizeTile(tile Tile, layer Layer) {
	newPalette := reducePalette(*tile.img, layer)

	for y := 0; y < tile.height; y++ {
		for x := 0; x < tile.width; x++ {
			pixelIndex := (tile.y+y)*tile.img.width + (tile.x + x)
			pixel := &tile.img.pixels[pixelIndex]
			if !pixel.HasBitPattern() { // has already been processed
				if layer.isLast { // last layer, all remaining pixels should be quantized against new palette
					quantizePixel(pixel, newPalette.palette)
					pixel.BitPattern = newPalette.bitpatterns[pixel.PaletteIndex]
				} else { // not the last layer, only process pixels that quantize to a bitpattern in the new palette
					quantizePixel(pixel, tile.img.palette)
					bitpattern, present := newPalette.bitpatterns[pixel.PaletteIndex]
					if present {
						pixel.BitPattern = bitpattern
					}
				}
			}
		}
	}
}

func Quantize(img IndexedImage) IndexedImage {
	result := img

	for _, layer := range img.spec.layers {
		// cut the image up according to layer specs
		tiles := getTiles(result, layer)

		// quantize the tiles
		for _, tile := range tiles {
			quantizeTile(tile, layer)
		}
	}
	return result
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func reducePalette(img IndexedImage, layer Layer) ReducedPalette {

	indexToCount := make(map[int]int)

	// TODO: does this make any difference?
	existingBitpatterns := make(map[int]int)

	// count nr of pixels for each quantized color
	for _, pixel := range img.pixels {

		// pixels that are already assigned a bitpattern don't count
		if !pixel.HasBitPattern() {
			quantizePixel(&pixel, img.palette)
			indexToCount[pixel.PaletteIndex]++
		} else {
			existingBitpatterns[pixel.PaletteIndex] = pixel.BitPattern
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

	newPalette := make(pixels.Palette)
	newBitpatterns := make(map[int]int)

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
