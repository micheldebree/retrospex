package conversion

import (
	"sort"

	"github.com/micheldebree/retrospex/internal/indexedimage"

	"github.com/micheldebree/retrospex/internal/pixels"
	"golang.org/x/exp/maps"
)

// PaletteDistance The distance from an RGB pixel to each key in a Palette
// In the order of the palette
type PaletteDistance map[int]float64

// QuantizePixel quantizes a pixel to the nearest color in the palette
func QuantizePixel(p *pixels.Pixel, pal pixels.Palette) {
	p.PaletteIndex = bestIndex(p, pal)
	// TODO: store quantization error for when we want to do error diffusion
	// dithering later
	p.QuantizationError = 0
}

// borrowed from golang color package
func bestIndex(p *pixels.Pixel, pal pixels.Palette) int {
	cr, cg, cb, ca := p.Color.RGBA()
	ret, bestSum := 0, uint32(1<<32-1)
	for i, v := range pal {
		vr, vg, vb, va := v.RGBA()
		sum := sqDiff(cr, vr) + sqDiff(cg, vg) + sqDiff(cb, vb) + sqDiff(ca, va)
		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}
	return ret
}

// borrowed from golang color package
func sqDiff(x, y uint32) uint32 {
	d := x - y
	return (d * d) >> 2
}

// Cut up image into regions for a particular layer
func getRegions(img indexedimage.IndexedImage, layer indexedimage.Layer) []Region {
	w, h := img.Width, img.Height

	nrCols, nrRows := w/layer.CellWidth, h/layer.CellHeight

	regions := make([]Region, nrCols*nrRows)

	for cy := range nrRows {
		for cx := range nrCols {

			regionIndex := cy*nrCols + cx

			regions[regionIndex] = Region{
				&img,
				cx * layer.CellWidth,
				cy * layer.CellHeight,
				layer.CellWidth,
				layer.CellHeight,
				make(map[int]int),
				make(map[int]int),
				layer.IsLast,
			}
			regions[regionIndex].initBitPatterns(layer)
		}
	}
	return regions
}

func quantizeRegion(region Region) {
	//region.assignColorToBitPattern(0b00, 0)
	assignBitPatterns(region)
	localPalette := region.getPalette()

	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.Pixels[pixelIndex]
			// only process pixels that don't have a bitpattern assigned yet
			if !pixel.HasBitPattern() {
				if region.isLastLayer {
					// last layer, all remaining pixels should be re-quantized
					// against the local palette
					QuantizePixel(pixel, localPalette)
				}
				region.AssignBitpatternToPixel(pixel)
			}
		}
	}
}

// quantize all the pixels in the image according to the image specs
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

	unassignedColors := getUnassignedColors(region)
	unmappedBitPatterns := region.getUnmappedBitpatterns()

	nrOfUnassignedColors := len(unassignedColors)
	nrOfUnmappedBitPatterns := len(unmappedBitPatterns)

	if nrOfUnassignedColors <= 0 || nrOfUnmappedBitPatterns <= 0 {
		return
	}

	index := 0
	done := false

	for !done {
		bitpattern := unmappedBitPatterns[index]
		region.assignColorToBitPattern(bitpattern, unassignedColors[index])
		index++
		done = index >= nrOfUnassignedColors || index >= nrOfUnmappedBitPatterns
	}

}

// Get colors that don't have a bitpattern assigned
// Sorted most occuring color first
func getUnassignedColors(region Region) []int {
	// maps color index to occurence count
	indexToCount := make(map[int]int)

	// count nr of pixels for each quantized color
	for y := 0; y < region.height; y++ {
		for x := 0; x < region.width; x++ {
			pixelIndex := region.coordsToIndex(x, y)
			pixel := &region.img.Pixels[pixelIndex]

			// pixels that are already assigned a bitpattern don't count
			if !pixel.HasBitPattern() {
				QuantizePixel(pixel, region.img.Palette)

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
