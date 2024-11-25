package indexedimage

import (
	"image"

	"github.com/micheldebree/retrospex/internal/pixels"
)

var RetrospecFactories = map[string]func(*image.Image) Retrospec{
	"koala":        makeKoalaSpec,
	"hires":        makeHiresSpec,
	"mixedcharset": makeMixedCharsetSpecSpec,
	"mccharset":    makeMCCharsetSpecSpec,
	"scccharset":   makeSCCCharsetSpecSpec,
	"mcibitmap":    makeMciBitmapSpec,
	"scsprites":    makeScSpritesSpec,
	"mcsprites":    makeMcSpritesSpec,
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

func makeMixedCharsetSpecSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00, 0x01, 0x10}, false}, // d021, d022, d023
			{4, 8, []int{0x11}, true},              // d800,x
		},
	}
}

func makeMCCharsetSpecSpec(img *image.Image) Retrospec {
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

func makeMciBitmapSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},            // d021
			{8, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		},
	}
}

func makeScSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0}, false},  // d021
			{24, 21, []int{1}, true}, // sprite color
		},
	}
}

func makeMcSpritesSpec(img *image.Image) Retrospec {
	w, h := pixels.GetDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},   // d021
			{12, 21, []int{0x01}, false}, // d025
			{12, 21, []int{0x10}, false}, // d027,x
			{12, 21, []int{0x11}, true},  // d026
		},
	}

}
