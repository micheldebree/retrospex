package conversion

import (
	"math"
	"sort"

	"github.com/micheldebree/retrospex/internal/indexedimage"

	"github.com/micheldebree/retrospex/internal/pixels"
	"golang.org/x/exp/maps"
)

// QuantizePixel quantizes a pixel to the nearest color in the palette
func QuantizePixel(p *pixels.Pixel, pal pixels.Palette, colorspace ColorspaceName) {
	p.PaletteIndex, p.QuantizationError = bestIndex(p, pal, colorspace)
}

type colorSpaceColor struct {
	c1, c2, c3 float64
}

func bestIndex(p *pixels.Pixel, pal pixels.Palette, colorspace ColorspaceName) (int, float64) {

	bestIndex, shortestDistance := indexedimage.UNKNOWN, math.MaxFloat64

	pColorSpace1, pColorSpace2, pColorSpace3 := toColorSpace(p.Color, colorspace)

	// Pre-convert palette if it's large, but here it's passed every time.
	// We can't easily cache it here without changing the signature.

	for i, v := range pal {
		vColorSpace1, vColorSpace2, vColorSpace3 := toColorSpace(v, colorspace)

		diff1 := (vColorSpace1 - pColorSpace1)
		diff2 := (vColorSpace2 - pColorSpace2)
		diff3 := (vColorSpace3 - pColorSpace3)

		distance := diff1*diff1 + diff2*diff2 + diff3*diff3

		if distance < shortestDistance {
			bestIndex, shortestDistance = i, distance
		}
	}

	return bestIndex, shortestDistance
}

func quantizeRegion(region indexedimage.Region, bitpatternToColor map[int]int, colorspace ColorspaceName, convertedPalette map[int]colorSpaceColor) {

	// initial mapping so we are able to force color to bitpattern assignment
	// beforehand
	for bitpattern, colorIndex := range bitpatternToColor {
		region.AssignColorToBitPattern(bitpattern, colorIndex)
	}

	assignBitPatterns(region, colorspace, convertedPalette)
	localPalette := region.GetPalette()

	// cache for quantized colors
	type colorKey struct {
		r, g, b, a uint32
	}
	colorCache := make(map[colorKey]int)

	for y := 0; y < region.Height; y++ {
		for x := 0; x < region.Width; x++ {
			pixel := region.GetPixel(x, y)
			// only process pixels that don't have a bitpattern assigned yet
			if !pixel.HasBitPattern() {
				if region.IsLastLayer {
					// last layer, all remaining pixels should be re-quantized
					// against the local palette
					r, g, b, a := pixel.Color.RGBA()
					key := colorKey{r, g, b, a}
					if paletteIndex, ok := colorCache[key]; ok {
						pixel.PaletteIndex = paletteIndex
					} else {
						bestIndexWithConvertedPalette(pixel, localPalette, colorspace, convertedPalette)
						colorCache[key] = pixel.PaletteIndex
					}
				}
				region.AssignBitpatternToPixel(pixel)
			}
		}
	}
}

// quantize all the pixels in the image according to the image specs
func Quantize(img indexedimage.IndexedImage, bitpatternToColor map[int]int, colorspace ColorspaceName) indexedimage.IndexedImage {
	result := img

	convertedPalette := make(map[int]colorSpaceColor, len(img.Palette))
	for i, v := range img.Palette {
		c1, c2, c3 := toColorSpace(v, colorspace)
		convertedPalette[i] = colorSpaceColor{c1, c2, c3}
	}

	for layerIndex := range img.Spec.Layers {

		// quantize the regions
		for _, region := range img.Regions[layerIndex] {
			quantizeRegion(region, bitpatternToColor, colorspace, convertedPalette)
		}
	}
	return result
}

// reduce a palette to maximum number of colors according to their
// quantized occurence in pixels. assign a bitpattern to each palette entry
// only considers pixels that don't have a bitpattern assigned yet
func bestIndexWithConvertedPalette(p *pixels.Pixel, pal pixels.Palette, colorspace ColorspaceName, convertedPalette map[int]colorSpaceColor) {
	bestIndex, shortestDistance := indexedimage.UNKNOWN, math.MaxFloat64

	pColorSpace1, pColorSpace2, pColorSpace3 := toColorSpace(p.Color, colorspace)

	for i := range pal {
		cv := convertedPalette[i]

		diff1 := (cv.c1 - pColorSpace1)
		diff2 := (cv.c2 - pColorSpace2)
		diff3 := (cv.c3 - pColorSpace3)

		distance := diff1*diff1 + diff2*diff2 + diff3*diff3

		if distance < shortestDistance {
			bestIndex, shortestDistance = i, distance
		}
	}

	p.PaletteIndex = bestIndex
	p.QuantizationError = shortestDistance
}

func assignBitPatterns(region indexedimage.Region, colorspace ColorspaceName, convertedPalette map[int]colorSpaceColor) {

	unassignedColors := getUnassignedColors(region, colorspace, convertedPalette)
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
func getUnassignedColors(region indexedimage.Region, colorspace ColorspaceName, convertedPalette map[int]colorSpaceColor) []int {
	// maps color index to occurence count
	indexToCount := make(map[int]int)

	// cache for quantized colors
	type colorKey struct {
		r, g, b, a uint32
	}
	colorCache := make(map[colorKey]int)

	// count nr of pixels for each quantized color
	for y := 0; y < region.Height; y++ {
		for x := 0; x < region.Width; x++ {
			pixel := region.GetPixel(x, y)

			// pixels that are already assigned a bitpattern don't count
			if !pixel.HasBitPattern() {
				r, g, b, a := pixel.Color.RGBA()
				key := colorKey{r, g, b, a}
				if paletteIndex, ok := colorCache[key]; ok {
					pixel.PaletteIndex = paletteIndex
				} else {
					bestIndexWithConvertedPalette(pixel, region.Img.Palette, colorspace, convertedPalette)
					colorCache[key] = pixel.PaletteIndex
				}

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
