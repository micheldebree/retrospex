package indexedimage

import (
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

const (
	KoalaType     RetrospecName = "koala"
	HiresType     RetrospecName = "hires"
	MCCharsetType RetrospecName = "mccharset"
	SCCharsetType RetrospecName = "sccharset"
	SCSpritesType RetrospecName = "scsprites"
	MCSpritesType RetrospecName = "mcsprites"
)

var BitsPerPixel = map[RetrospecName]int{
	KoalaType:     2,
	HiresType:     1,
	MCCharsetType: 2,
	SCCharsetType: 1,
	SCSpritesType: 1,
	MCSpritesType: 2,
}

var RetrospecFactories = map[RetrospecName]func(*image.Image) Retrospec{
	KoalaType:     makeKoalaSpec,
	HiresType:     makeHiresSpec,
	MCCharsetType: makeMCCharsetSpec,
	SCCharsetType: makeSCCCharsetSpec,
	SCSpritesType: makeSCSpritesSpec,
	MCSpritesType: makeMCSpritesSpec,
}

var RetrospecTemplateFactories = map[RetrospecName]func() RetrospecTemplate{
	KoalaType:     makeKoalaTemplate,
	HiresType:     makeHiresTemplate,
	MCCharsetType: makeMCCharsetTemplate,
	SCCharsetType: makeSCCCharsetTemplate,
	SCSpritesType: makeSCSpritesTemplate,
	MCSpritesType: makeMCSpritesTemplate,
}

func MakeSpec(specName RetrospecName, img *image.Image) Retrospec {
	template := MakeTemplate(specName)
	w, h := pixels.GetDimensions(img)

	// Convert template layers to final layers with proper dimensions
	layers := make([]Layer, len(template.Layers))
	for i, layerTemplate := range template.Layers {
		if layerTemplate.UseImageDimensions {
			layers[i] = Layer{w, h, layerTemplate.Bitpatterns, layerTemplate.IsLast}
		} else {
			layers[i] = Layer{layerTemplate.CellWidth, layerTemplate.CellHeight, layerTemplate.Bitpatterns, layerTemplate.IsLast}
		}
	}

	return Retrospec{template.Name, layers, template.BitsPerPixel}
}

// GetBitPatternSize returns the hardcoded bitpattern size for a given mode
func GetBitPatternSize(modeName RetrospecName) int {
	size, isPresent := BitsPerPixel[modeName]
	if !isPresent {
		panic("Unknown mode")
	}
	return size
}

// MakeTemplate creates a template spec for the given mode (for resize operations)
func MakeTemplate(specName RetrospecName) RetrospecTemplate {
	factory, isPresent := RetrospecTemplateFactories[specName]
	if !isPresent {
		panic("Unknown mode")
	}
	return factory()
}

// +----+----+----+----+----+----+----+----+----+----+----+----+
// | 11 | 10 |  9 |  8 |  7 |  6 |  5 |  4 |  3 |  2 |  1 |  0 |
// +----+----+----+----+----+----+----+----+----+----+----+----+
// |     Color of      |     Color of      |     Color of      |
// |    "11 pixels"    |    "01" pixels    |    "10" pixels    |
// +-------------------+-------------------+-------------------+

func makeKoalaSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{KoalaType,
		[]Layer{
			{w, h, []int{0b00}, false},            // d021
			{4, 8, []int{0b01, 0b10, 0b11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		}, 2,
	}
}

func makeHiresSpec(_ *image.Image) Retrospec {
	return Retrospec{HiresType,
		[]Layer{
			{8, 8, []int{0b00, 0b01}, true}, // 0400,x (lower nibble), 0400,x (upper nibble)
		}, 1,
	}
}

func makeMCCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{MCCharsetType,
		[]Layer{
			{w, h, []int{0b00, 0b01, 0b10}, false}, // d021, d022, d023
			{4, 8, []int{0b11}, true},              // d800,x
		}, 2,
	}
}

func makeSCCCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{SCCharsetType,
		[]Layer{
			{w, h, []int{0b00}, false}, // d021
			{8, 8, []int{0b01}, true},  // d800,x
		}, 1,
	}
}

// TODO: does png2prg only support one color?
func makeSCSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{SCSpritesType,
		[]Layer{
			{w, h, []int{0b00, 0b01}, true}, // d021, d027,x
		}, 1,
	}
}

// TODO: does png2prg only support 4 colors?
func makeMCSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{MCSpritesType,
		[]Layer{
			{w, h, []int{0b00, 0b01, 0b10, 0b11}, false}, // d021, d025, d027,x, d026
		}, 2,
	}
}

// Template factory functions
func makeKoalaTemplate() RetrospecTemplate {
	return RetrospecTemplate{KoalaType,
		[]LayerTemplate{
			{true, 0, 0, []int{0b00}, false},             // d021
			{false, 4, 8, []int{0b01, 0b10, 0b11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		}, 2,
	}
}

func makeHiresTemplate() RetrospecTemplate {
	return RetrospecTemplate{HiresType,
		[]LayerTemplate{
			{false, 8, 8, []int{0b00, 0b01}, true}, // 0400,x (lower nibble), 0400,x (upper nibble)
		}, 1,
	}
}

func makeMCCharsetTemplate() RetrospecTemplate {
	return RetrospecTemplate{MCCharsetType,
		[]LayerTemplate{
			{true, 0, 0, []int{0b00, 0b01, 0b10}, false}, // d021, d022, d023
			{false, 4, 8, []int{0b11}, true},             // d800,x
		}, 2,
	}
}

func makeSCCCharsetTemplate() RetrospecTemplate {
	return RetrospecTemplate{SCCharsetType,
		[]LayerTemplate{
			{true, 0, 0, []int{0b00}, false}, // d021
			{false, 8, 8, []int{0b01}, true}, // d800,x
		}, 1,
	}
}

func makeSCSpritesTemplate() RetrospecTemplate {
	return RetrospecTemplate{SCSpritesType,
		[]LayerTemplate{
			{true, 0, 0, []int{0b00, 0b01}, true}, // d021, d027,x
		}, 1,
	}
}

func makeMCSpritesTemplate() RetrospecTemplate {
	return RetrospecTemplate{MCSpritesType,
		[]LayerTemplate{
			{true, 0, 0, []int{0b00, 0b01, 0b10, 0b11}, false}, // d021, d025, d027,x, d026
		}, 2,
	}
}
