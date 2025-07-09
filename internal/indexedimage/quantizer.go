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
	distances := distances(p.Color, pal)
	bestIndex, smallestDistance := bestPixelIndex(distances)
	p.PaletteIndex = bestIndex
	p.QuantizationError = smallestDistance
}

// LayerRegion represents a region of an IndexedImage with its reduced palette of colors
// It belongs to a layer and maps the bitpatterns in that layer to this palette
type LayerRegion struct {
	img           *IndexedImage
	layer         *Layer
	x, y          int
	width, height int
	palette       pixels.Palette // a 'sparse' palette, meaning not all palette indexes from the original palette are present
	bitpatterns   map[int]int    // maps each palette index to a bit pattern
}

func (region *LayerRegion) coordsToIndex(x, y int) int {
	return (region.y+y)*region.img.width + (region.x + x)
}

// Cut up image into regions for a particular layer
func getLayerRegions(img IndexedImage, layer Layer) []LayerRegion {
	w, h := img.width, img.height

	nrCols, nrRows := w/layer.cellWidth, h/layer.cellHeight

	regions := make([]LayerRegion, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {
			regions[cy*nrCols+cx] = LayerRegion{
				&img,
				&layer,
				cx * layer.cellWidth,
				cy * layer.cellHeight,
				layer.cellWidth,
				layer.cellHeight,
				pixels.Palette{},
				make(map[int]int),
			}
		}
	}
	return regions
}

func quantizeLayerRegion(region LayerRegion) {
	createPaletteForBitpatterns(region)

	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.pixels[pixelIndex]
			if !pixel.HasBitPattern() { // has already been processed
				if region.layer.isLast { // last layer, all remaining pixels should be quantized against the region's palette
					quantizePixel(pixel, region.palette)
					pixel.BitPattern = region.bitpatterns[pixel.PaletteIndex]
				} else { // not the last layer, only process pixels that quantize to a bitpattern in the new palette
					quantizePixel(pixel, region.img.palette)
					bitpattern, present := region.bitpatterns[pixel.PaletteIndex]
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
		regions := getLayerRegions(result, layer)

		// quantize the regions
		for _, region := range regions {
			quantizeLayerRegion(region)
		}
	}
	return result
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func createPaletteForBitpatterns(region LayerRegion) {

	indexToCount := make(map[int]int)

	// count nr of pixels for each quantized color
	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.pixels[pixelIndex]

			// pixels that are already assigned a bitpattern don't count
			if !pixel.HasBitPattern() {
				quantizePixel(pixel, region.img.palette)
				indexToCount[pixel.PaletteIndex]++
			}
		}
	}

	// sort in reverse order of count values
	keys := maps.Keys(indexToCount)
	sort.SliceStable(keys, func(i, j int) bool {
		return indexToCount[keys[i]] > indexToCount[keys[j]]
	})

	// only keep top n
	maxColors := len(region.layer.bitpatterns)
	if maxColors < len(keys) {
		keys = keys[0:maxColors]
	}

	// assign bitpatterns
	i := 0
	for _, key := range keys {
		region.palette[key] = region.img.palette[key]
		region.bitpatterns[key] = region.layer.bitpatterns[i]
		i++
	}

}
