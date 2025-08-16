package indexedimage

import (
	"fmt"

	"github.com/micheldebree/retrospex/internal/pixels"
)

const UNKNOWN = -1

// Region represents a region of an IndexedImage
// Regions can overlap
// Within one region, there is a limited palette, assigned to the bitpatterns in that region
type Region struct {
	Img               *IndexedImage
	x, y              int
	Width, Height     int
	bitpatternToColor map[int]int // maps each bit pattern to a palette index
	colorToBitpattern map[int]int // reverse lookup for bitpatternToColor (optimization)
	IsLastLayer       bool        // is this region in the last layer?
}

func (region *Region) initBitPatterns(layer Layer) {
	for _, bitPattern := range layer.Bitpatterns {
		region.bitpatternToColor[bitPattern] = UNKNOWN
	}
}

func (region *Region) bitPatternIsAssigned(colorIndex int) bool {
	colorIndex, present := region.bitpatternToColor[colorIndex]
	return present && colorIndex != UNKNOWN
}

// associates a bitpattern with an index in the color palette
func (region *Region) AssignColorToBitPattern(bitPattern int, paletteIndex int) {
	if region.bitPatternIsAssigned(bitPattern) {
		panic(fmt.Sprintf("cannot assign color %d to bit pattern %d, a color is already assigned for region %v", paletteIndex, bitPattern, region))
	}
	region.bitpatternToColor[bitPattern] = paletteIndex
	region.colorToBitpattern[paletteIndex] = bitPattern
}

// if a bitpattern is associated with the pixel's quantized color, assign the bitpattern to the pixel
func (region *Region) AssignBitpatternToPixel(pixel *pixels.Pixel) bool {
	pixel.AssertQuantized()

	if bitpattern, present := region.colorToBitpattern[pixel.PaletteIndex]; present {
		pixel.BitPattern = bitpattern
		return true
	}
	return false
}

func (region *Region) GetPixel(x, y int) *pixels.Pixel {
	pixelIndex := (region.y+y)*region.Img.Width + (region.x + x)
	return &region.Img.Pixels[pixelIndex]
}

func (region *Region) GetUnmappedBitpatterns() []int {
	result := make([]int, 0)
	for bitpattern, paletteIndex := range region.bitpatternToColor {
		if paletteIndex == UNKNOWN {
			result = append(result, bitpattern)
		}
	}
	return result
}

// create a pallete containing only the colors assigned to a bitpattern
func (region *Region) GetPalette() pixels.Palette {

	result := make(pixels.Palette)
	for _, paletteIndex := range region.bitpatternToColor {
		if paletteIndex >= 0 {
			result[paletteIndex] = region.Img.Palette[paletteIndex]
		}
	}
	return result
}
