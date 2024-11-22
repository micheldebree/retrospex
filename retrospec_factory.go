package main

import "image"

var RetrospecFactories = map[string]func(*image.Image) Retrospec{
	"koala":        makeKoalaSpec,
	"hires":        makeHiresSpec,
	"mixedcharset": makeMixedCharsetSpecSpec,
	"mccharset":    makeMCCharsetSpecSpec,
	"scccharset":   makeSCCCharsetSpecSpec,
	"mcibitmap":    makeMciBitmapSpec,
}

func makeSpec(specName string, img *image.Image) Retrospec {

	factory, isPresent := RetrospecFactories[specName]
	if !isPresent {
		panic("Unknown mode")
	}
	return factory(img)
}

func makeKoalaSpec(img *image.Image) Retrospec {
	w, h := getDimensions(img)
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
	w, h := getDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00, 0x01, 0x10}, false}, // d021, d022, d023
			{4, 8, []int{0x11}, true},              // d800,x
		},
	}
}

func makeMCCharsetSpecSpec(img *image.Image) Retrospec {
	w, h := getDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00, 0x01, 0x10, 0x11}, true}, // d021, d022, d023, d800...
		},
	}
}

func makeSCCCharsetSpecSpec(img *image.Image) Retrospec {
	w, h := getDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0}, false}, // d021
			{8, 8, []int{1}, true},  // d800,x
		},
	}
}

func makeMciBitmapSpec(img *image.Image) Retrospec {
	w, h := getDimensions(img)
	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},            // d021
			{8, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		},
	}
}
