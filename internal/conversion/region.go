package conversion

import (
	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

// Region represents a region of an IndexedImage
// Regions can overlap
// Within one region, there is a limited palette, assigned to the bitpatterns in that region
type Region struct {
	img               *indexedimage.IndexedImage
	x, y              int
	width, height     int
	bitpatternToColor map[int]int // maps each bit pattern to a palette index
	colorToBitpattern map[int]int // maps each bit pattern to a palette index
	isLastLayer       bool        // is this region in the last layer?
}

func (region *Region) initBitPatterns(layer indexedimage.Layer) {
	for _, bitPattern := range layer.Bitpatterns {
		region.bitpatternToColor[bitPattern] = -1
	}
}

func (region *Region) bitPatternIsAssigned(bitPattern int) bool {
	return region.bitpatternToColor[bitPattern] != -1
}

func (region *Region) assignColorToBitPattern(bitPattern int, paletteIndex int) {
	if region.bitPatternIsAssigned(bitPattern) {
		panic("cannot assign color to bit pattern, a color is already assigned")
	}
	region.bitpatternToColor[bitPattern] = paletteIndex
	region.colorToBitpattern[paletteIndex] = bitPattern
}

func (region *Region) coordsToIndex(x, y int) int {
	return (region.y+y)*region.img.Width + (region.x + x)
}

// get the first bitpattern that is not mapped to a pallette index
// TODO: the map is unordered so unpredictable
func (region *Region) getUnmappedBitPattern() int {
	for bitpattern, paletteIndex := range region.bitpatternToColor {
		if paletteIndex < 0 {
			return bitpattern
		}
	}
	panic("No more bitpatterns to assign")
}

// create a pallete containing only the colors assigned to a bitpattern
func (region *Region) getPalette() pixels.Palette {

	result := make(pixels.Palette)
	for _, paletteIndex := range region.bitpatternToColor {
		if paletteIndex >= 0 {
			result[paletteIndex] = region.img.Palette[paletteIndex]
		}
	}
	return result
}
