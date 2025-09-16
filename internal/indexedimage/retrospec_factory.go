package indexedimage

import (
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

const (
	KoalaType        RetrospecName = "koala"
	HiresType        RetrospecName = "hires"
	MixedCharsetType RetrospecName = "mixedcharset"
	MCCharsetType    RetrospecName = "mccharset"
	SCCharsetType    RetrospecName = "scccharset"
	MCIBitmapType    RetrospecName = "mcibitmap"
	SCSpritesType    RetrospecName = "scsprites"
	MCSpritesType    RetrospecName = "mcsprites"
)

var RetrospecFactories = map[RetrospecName]func(*image.Image) Retrospec{
	KoalaType:        makeKoalaSpec,
	HiresType:        makeHiresSpec,
	MixedCharsetType: makeMixedCharsetSpec,
	MCCharsetType:    makeMCCharsetSpec,
	SCCharsetType:    makeSCCCharsetSpec,
	MCIBitmapType:    makeMCiBitmapSpec,
	SCSpritesType:    makeSCSpritesSpec,
	MCSpritesType:    makeMCSpritesSpec,
}

func MakeSpec(specName RetrospecName, img *image.Image) Retrospec {
	factory, isPresent := RetrospecFactories[specName]
	if !isPresent {
		panic("Unknown mode")
	}
	return factory(img)
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

func makeMixedCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{MixedCharsetType,
		[]Layer{
			{w, h, []int{0b00, 0b01, 0b10}, false}, // d021, d022, d023
			{4, 8, []int{0b11}, true},              // d800,x
		}, 2,
	}
}

func makeMCCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{MCCharsetType,
		[]Layer{
			{w, h, []int{0b00, 0b01, 0b10, 0b11}, true}, // d021, d022, d023, d800...
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

func makeMCiBitmapSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{MCIBitmapType,
		[]Layer{
			{w, h, []int{0b00}, false},            // d021
			{8, 8, []int{0b01, 0b10, 0b11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		}, 2,
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
