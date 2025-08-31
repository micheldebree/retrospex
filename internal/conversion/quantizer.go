package conversion

import (
	"math"
	"sort"

	"github.com/micheldebree/retrospex/internal/indexedimage"

	"github.com/micheldebree/retrospex/internal/pixels"
	"golang.org/x/exp/maps"
)

// PaletteDistance The distance from an RGB pixel to each key in a Palette
// In the order of the palette
type PaletteDistance map[int]float64

// QuantizePixel quantizes a pixel to the nearest color in the palette
func QuantizePixel(p *pixels.Pixel, pal pixels.Palette, colorspace ColorspaceName) {
	p.PaletteIndex, p.QuantizationError = convertertedBestIndex(p, pal, colorspace)
}

func convertertedBestIndex(p *pixels.Pixel, pal pixels.Palette, colorspace ColorspaceName) (int, float64) {

	bestIndex, shortestDistance := indexedimage.UNKNOWN, math.MaxFloat64

	for i, v := range pal {
		distance := euclidianDistance(p.Color, v, colorspace)

		if distance < shortestDistance {
			bestIndex, shortestDistance = i, distance
		}
	}

	return bestIndex, shortestDistance
}

func quantizeRegion(region indexedimage.Region, bitpatternToColor map[int]int, colorspace ColorspaceName) {

	// initial mapping so we are able to force color to bitpattern assignment
	// beforehand
	for bitpattern, colorIndex := range bitpatternToColor {
		region.AssignColorToBitPattern(bitpattern, colorIndex)
	}

	assignBitPatterns(region, colorspace)
	localPalette := region.GetPalette()

	for y := 0; y < region.Height; y++ {
		for x := 0; x < region.Width; x++ {
			pixel := region.GetPixel(x, y)
			// only process pixels that don't have a bitpattern assigned yet
			if !pixel.HasBitPattern() {
				if region.IsLastLayer {
					// last layer, all remaining pixels should be re-quantized
					// against the local palette
					QuantizePixel(pixel, localPalette, colorspace)
				}
				region.AssignBitpatternToPixel(pixel)
			}
		}
	}
}

// quantize all the pixels in the image according to the image specs
func Quantize(img indexedimage.IndexedImage, bitpatternToColor map[int]int, colorspace ColorspaceName) indexedimage.IndexedImage {
	result := img

	for layerIndex := range img.Spec.Layers {

		// quantize the regions
		for _, region := range img.Regions[layerIndex] {
			quantizeRegion(region, bitpatternToColor, colorspace)
		}
	}
	return result
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func assignBitPatterns(region indexedimage.Region, colorspace ColorspaceName) {

	unassignedColors := getUnassignedColors(region, colorspace)
	unmappedBitPatterns := region.GetUnmappedBitpatterns()

	nrOfUnassignedColors := len(unassignedColors)
	nrOfUnmappedBitPatterns := len(unmappedBitPatterns)

	if nrOfUnassignedColors <= 0 || nrOfUnmappedBitPatterns <= 0 {
		return
	}

	index := 0
	done := false

	for !done {
		bitpattern := unmappedBitPatterns[index]
		region.AssignColorToBitPattern(bitpattern, unassignedColors[index])
		index++
		done = index >= nrOfUnassignedColors || index >= nrOfUnmappedBitPatterns
	}

}

// Get colors that don't have a bitpattern assigned
// Sorted most occuring color first
func getUnassignedColors(region indexedimage.Region, colorspace ColorspaceName) []int {
	// maps color index to occurence count
	indexToCount := make(map[int]int)

	// count nr of pixels for each quantized color
	for y := 0; y < region.Height; y++ {
		for x := 0; x < region.Width; x++ {
			pixel := region.GetPixel(x, y)

			// pixels that are already assigned a bitpattern don't count
			if !pixel.HasBitPattern() {
				QuantizePixel(pixel, region.Img.Palette, colorspace)

				// there could already be colors associated with bitpatterns
				// initially. If so, use it and don't count as unassigned
				if success := region.AssignBitpatternToPixel(pixel); !success {
					indexToCount[pixel.PaletteIndex]++
				}
			}
		}
	}

	// sort in reverse order of count values
	unassignedColors := maps.Keys(indexToCount)
	sort.SliceStable(unassignedColors, func(i, j int) bool {
		return indexToCount[unassignedColors[i]] > indexToCount[unassignedColors[j]]
	})
	return unassignedColors
}
