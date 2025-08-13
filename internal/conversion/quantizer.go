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
					QuantizePixel(pixel, localPalette)
					pixel.BitPattern = region.colorToBitpattern[pixel.PaletteIndex]
				} else { // not the last layer, only process pixels that quantize to a bitpattern in the new palette
					QuantizePixel(pixel, region.img.Palette)
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
				QuantizePixel(pixel, region.img.Palette)
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
