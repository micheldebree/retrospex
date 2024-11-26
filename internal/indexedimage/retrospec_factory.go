package indexedimage

import (
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

var RetrospecFactories = map[string]func(*image.Image) Retrospec{
	"koala":        makeKoalaSpec,
	"hires":        makeHiresSpec,
	"mixedcharset": makeMixedCharsetSpec,
	"mccharset":    makeMCCharsetSpec,
	"scccharset":   makeSCCCharsetSpecSpec,
	"mcibitmap":    makeMCiBitmapSpec,
	"scsprites":    makeSCSpritesSpec,
	"mcsprites":    makeMCSpritesSpec,
}

func MakeSpec(specName string, img *image.Image) Retrospec {
	factory, isPresent := RetrospecFactories[specName]
	if !isPresent {
		panic("Unknown mode")
	}
	return factory(img)
}

func makeKoalaSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},            // d021
			{4, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		},
	}
}

func makeHiresSpec(img *image.Image) Retrospec {
	return Retrospec{
		[]Layer{
			{8, 8, []int{0, 1}, true}, // 0400,x (lower nibble), 0400,x (upper nibble)
		},
	}
}

func makeMixedCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00, 0x01, 0x10}, false}, // d021, d022, d023
			{4, 8, []int{0x11}, true},              // d800,x
		},
	}
}

func makeMCCharsetSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00, 0x01, 0x10, 0x11}, true}, // d021, d022, d023, d800...
		},
	}
}

func makeSCCCharsetSpecSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0}, false}, // d021
			{8, 8, []int{1}, true},  // d800,x
		},
	}
}

func makeMCiBitmapSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},            // d021
			{8, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		},
	}
}

// TODO: does png2prg only support one color?
func makeSCSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0}, false}, // d021
			{w, h, []int{1}, true},  // sprite color
		},
	}
}

// TODO: does png2prg only support 4 colors?
func makeMCSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false}, // d021
			{w, h, []int{0x01}, false}, // d025
			{w, h, []int{0x10}, false}, // d027,x
			{w, h, []int{0x11}, true},  // d026
		},
	}
}
