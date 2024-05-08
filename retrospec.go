package main

type Layer struct {
	cellWidth, cellHeight int
	bitpatterns           []int
	// TODO: this can be determined while processing
	isLast bool // the last layer should quantize all remaining pixels
}

type Retrospec struct {
	width      int
	height     int
	pixelWidth int
	layers     []Layer
}

func (spec Retrospec) displayWidth() int {
	return spec.width * spec.pixelWidth
}

func (spec Retrospec) displayHeight() int {
	return spec.height
}

var C64Specs = map[string]Retrospec{
	"koala":        KoalaSpec,
	"hires":        HiresSpec,
	"mixedcharset": MixedCharsetSpec,
	"mccharset":    MCCharsetSpec,
	"sccharset":    SCCharsetSpec,
	"mcibitmap":    MciBitmapSpec,
}

var KoalaSpec = Retrospec{160, 200, 2,
	[]Layer{
		{160, 200, []int{0x00}, false},        // d021
		{4, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
	},
}

var HiresSpec = Retrospec{320, 200, 1,
	[]Layer{
		{8, 8, []int{0, 1}, true}, // 0400,x (lower nibble), 0400,x (upper nibble)
	},
}

var MixedCharsetSpec = Retrospec{160, 200, 2,
	[]Layer{
		{160, 200, []int{0x00, 0x01, 0x10}, false}, // d021, d022, d023
		{4, 8, []int{0x11}, true},                  // d800,x
	},
}

var SCCharsetSpec = Retrospec{320, 200, 1,
	[]Layer{
		{320, 200, []int{0}, false}, // d021
		{8, 8, []int{1}, true},      // d800,x
	},
}

var Hires2Colors = Retrospec{320, 200, 1,
	[]Layer{
		{320, 200, []int{0, 1}, true}, // 0400,x (lower), 0400,x (higher)
	},
}

var MCCharsetSpec = Retrospec{160, 200, 2,
	[]Layer{
		{160, 200, []int{0x00, 0x01, 0x10, 0x11}, true}, // d021, d022, d023, d800...
	},
}

var MciBitmapSpec = Retrospec{320, 200, 1,
	[]Layer{
		{320, 200, []int{0x00}, false},        // d021
		{8, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
	},
}
