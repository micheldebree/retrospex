package conversion

import (
	"image/color"
	"math"
	"sort"

	"github.com/micheldebree/retrospex/internal/indexedimage"

	"github.com/micheldebree/retrospex/internal/pixels"
	"golang.org/x/exp/maps"
)

// PaletteDistance The distance from an RGB pixel to each key in a Palette
// In the order of the palette
type PaletteDistance map[int]float64

func distance(color1, color2 color.Color) float64 {
	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()

	dr := float64(r1 - r2)
	dg := float64(g1 - g2)
	db := float64(b1 - b2)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// distances Distance from a pixel to each color in a palette
func distances(aColor color.Color, palette pixels.Palette) PaletteDistance {
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

// Cut up image into regions for a particular layer
func getRegions(img indexedimage.IndexedImage, layer indexedimage.Layer) []Region {
	w, h := img.Width, img.Height

	nrCols, nrRows := w/layer.CellWidth, h/layer.CellHeight

	regions := make([]Region, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {
			regions[cy*nrCols+cx] = Region{
				&img,
				cx * layer.CellWidth,
				cy * layer.CellHeight,
				layer.CellWidth,
				layer.CellHeight,
				make(map[int]int),
				make(map[int]int),
				layer.IsLast,
			}

			// initialize bitpatterns, unmapped
			for _, bitpattern := range layer.Bitpatterns {
				regions[cy*nrCols+cx].addMapping(bitpattern, -1)
			}

		}
	}
	return regions
}

func quantizeRegion(region Region) {
	assignBitPatterns(region)
	localPalette := region.getPalette()

	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.Pixels[pixelIndex]
			if !pixel.HasBitPattern() { // has already been processed
				if region.isLastLayer { // last layer, all remaining pixels should be quantized against the region's palette
					quantizePixel(pixel, localPalette)
					pixel.BitPattern = region.colorToBitpattern[pixel.PaletteIndex]
				} else { // not the last layer, only process pixels that quantize to a bitpattern in the new palette
					quantizePixel(pixel, region.img.Palette)
					bitpattern, present := region.colorToBitpattern[pixel.PaletteIndex]
					if present {
						pixel.BitPattern = bitpattern
					}
				}
			}
		}
	}
}

func Quantize(img indexedimage.IndexedImage) indexedimage.IndexedImage {
	result := img

	for _, layer := range img.Spec.Layers {
		// cut the image up according to layer specs
		regions := getRegions(result, layer)

		// quantize the regions
		for _, region := range regions {
			quantizeRegion(region)
		}
	}
	return result
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func assignBitPatterns(region Region) {

	// maps color index to occurence count
	indexToCount := make(map[int]int)

	// count nr of pixels for each quantized color
	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.Pixels[pixelIndex]

			// pixels that are already assigned a bitpattern don't count
			if !pixel.HasBitPattern() {
				quantizePixel(pixel, region.img.Palette)
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
	maxColors := len(region.bitpatternToColor)
	if maxColors < len(keys) {
		keys = keys[0:maxColors]
	}

	// assign bitpatterns
	for _, key := range keys {
		bitPattern := region.getUnmappedBitPattern()
		region.addMapping(bitPattern, key)
	}

}
