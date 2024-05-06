package main

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type Palette map[int]colorful.Color

var C64Palettes = map[string]Palette{
	"colodore":          Colodore,
	"pepto":             Pepto,
	"vice":              Vice,
	"vice_old_lum":      ViceOldLum,
	"vice_371_internal": Vice371Internal,
	"pantaloon":         Pantaloon,
	"archmage":          Archmage,
	"electric":          Electric,
	"ste":               Ste,
	"perplex_1":         Perplex,
	"perplex_2":         Perplex2,
	"lenyn":             Lenyn,
	// TODO: PALette
}

var Vice = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xbc, 0x52, 0x41, 0xff},
		color.RGBA{0x8f, 0xef, 0xfb, 0xff},
		color.RGBA{0xb9, 0x56, 0xeb, 0xff},
		color.RGBA{0x7e, 0xdb, 0x40, 0xff},
		color.RGBA{0x55, 0x3f, 0xe4, 0xff},
		color.RGBA{0xff, 0xff, 0x77, 0xff},
		color.RGBA{0xc1, 0x7b, 0x1d, 0xff},
		color.RGBA{0x82, 0x63, 0x00, 0xff},
		color.RGBA{0xf4, 0x94, 0x86, 0xff},
		color.RGBA{0x72, 0x72, 0x72, 0xff},
		color.RGBA{0xa4, 0xa4, 0xa4, 0xff},
		color.RGBA{0xcd, 0xff, 0x98, 0xff},
		color.RGBA{0x9e, 0x8d, 0xff, 0xff},
		color.RGBA{0xd5, 0xd5, 0xd5, 0xff},
	})

var ViceOldLum = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xa9, 0x38, 0x26, 0xff},
		color.RGBA{0xae, 0xff, 0xff, 0xff},
		color.RGBA{0xdf, 0x82, 0xff, 0xff},
		color.RGBA{0x7e, 0xdb, 0x40, 0xff},
		color.RGBA{0x55, 0x3f, 0xe4, 0xff},
		color.RGBA{0xf7, 0xff, 0x6d, 0xff},
		color.RGBA{0xe7, 0xa4, 0x53, 0xff},
		color.RGBA{0x82, 0x63, 0x00, 0xff},
		color.RGBA{0xf4, 0x94, 0x86, 0xff},
		color.RGBA{0x5c, 0x5c, 0x5c, 0xff},
		color.RGBA{0xb0, 0xb0, 0xb0, 0xff},
		color.RGBA{0xc4, 0xff, 0x8f, 0xff},
		color.RGBA{0xaa, 0x99, 0xff, 0xff},
		color.RGBA{0xf2, 0xf2, 0xf2, 0xff},
	})

var Pepto = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x68, 0x37, 0x2b, 0xff},
		color.RGBA{0x70, 0xa4, 0xb2, 0xff},
		color.RGBA{0x6f, 0x3d, 0x86, 0xff},
		color.RGBA{0x58, 0x8d, 0x43, 0xff},
		color.RGBA{0x35, 0x28, 0x79, 0xff},
		color.RGBA{0xb8, 0xc7, 0x6f, 0xff},
		color.RGBA{0x6f, 0x4f, 0x25, 0xff},
		color.RGBA{0x43, 0x39, 0x00, 0xff},
		color.RGBA{0x9a, 0x67, 0x59, 0xff},
		color.RGBA{0x44, 0x44, 0x44, 0xff},
		color.RGBA{0x6c, 0x6c, 0x6c, 0xff},
		color.RGBA{0x9a, 0xd2, 0x84, 0xff},
		color.RGBA{0x6c, 0x5e, 0xb5, 0xff},
		color.RGBA{0x95, 0x95, 0x95, 0xff},
	})

var Pantaloon = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x68, 0x37, 0x2b, 0xff},
		color.RGBA{0x83, 0xf0, 0xdc, 0xff},
		color.RGBA{0x6f, 0x3d, 0x86, 0xff},
		color.RGBA{0x59, 0xcd, 0x36, 0xff},
		color.RGBA{0x41, 0x37, 0xcd, 0xff},
		color.RGBA{0xb8, 0xc7, 0x6f, 0xff},
		color.RGBA{0xd1, 0x7f, 0x30, 0xff},
		color.RGBA{0x43, 0x39, 0x00, 0xff},
		color.RGBA{0x9a, 0x67, 0x59, 0xff},
		color.RGBA{0x5b, 0x5b, 0x5b, 0xff},
		color.RGBA{0x8e, 0x8e, 0x8e, 0xff},
		color.RGBA{0x9d, 0xff, 0x9d, 0xff},
		color.RGBA{0x75, 0xa1, 0xec, 0xff},
		color.RGBA{0xc1, 0xc1, 0xc1, 0xff},
	})

var Archmage = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x89, 0x40, 0x37, 0xff},
		color.RGBA{0x7a, 0xbf, 0xc7, 0xff},
		color.RGBA{0x8a, 0x46, 0xae, 0xff},
		color.RGBA{0x68, 0xa9, 0x41, 0xff},
		color.RGBA{0x3e, 0x31, 0xa2, 0xff},
		color.RGBA{0xd0, 0xdc, 0x71, 0xff},
		color.RGBA{0x90, 0x5f, 0x25, 0xff},
		color.RGBA{0x5c, 0x47, 0x00, 0xff},
		color.RGBA{0xbb, 0x77, 0x6d, 0xff},
		color.RGBA{0x55, 0x55, 0x55, 0xff},
		color.RGBA{0x80, 0x80, 0x80, 0xff},
		color.RGBA{0xac, 0xea, 0x88, 0xff},
		color.RGBA{0x7c, 0x70, 0xda, 0xff},
		color.RGBA{0xab, 0xab, 0xab, 0xff},
	})

var Electric = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x8b, 0x1f, 0x00, 0xff},
		color.RGBA{0x6f, 0xdf, 0xb7, 0xff},
		color.RGBA{0xa7, 0x3b, 0x9f, 0xff},
		color.RGBA{0x4a, 0xb5, 0x10, 0xff},
		color.RGBA{0x08, 0x00, 0x94, 0xff},
		color.RGBA{0xf3, 0xeb, 0x5b, 0xff},
		color.RGBA{0xa5, 0x42, 0x00, 0xff},
		color.RGBA{0x63, 0x29, 0x18, 0xff},
		color.RGBA{0xcb, 0x7b, 0x6f, 0xff},
		color.RGBA{0x45, 0x44, 0x44, 0xff},
		color.RGBA{0x9f, 0x9f, 0x9f, 0xff},
		color.RGBA{0x94, 0xff, 0x94, 0xff},
		color.RGBA{0x4a, 0x94, 0xd6, 0xff},
		color.RGBA{0xbd, 0xbd, 0xbd, 0xff},
	})

var Ste = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xc8, 0x35, 0x35, 0xff},
		color.RGBA{0x83, 0xf0, 0xdc, 0xff},
		color.RGBA{0xcc, 0x59, 0xc6, 0xff},
		color.RGBA{0x59, 0xcd, 0x36, 0xff},
		color.RGBA{0x41, 0x37, 0xcd, 0xff},
		color.RGBA{0xf7, 0xee, 0x59, 0xff},
		color.RGBA{0xd1, 0x7f, 0x30, 0xff},
		color.RGBA{0x91, 0x5f, 0x33, 0xff},
		color.RGBA{0xf9, 0x9b, 0x97, 0xff},
		color.RGBA{0x5b, 0x5b, 0x5b, 0xff},
		color.RGBA{0x8e, 0x8e, 0x8e, 0xff},
		color.RGBA{0x9d, 0xff, 0x9d, 0xff},
		color.RGBA{0x75, 0xa1, 0xec, 0xff},
		color.RGBA{0xc1, 0xc1, 0xc1, 0xff},
	})

var Perplex = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xfd, 0xfe, 0xfc, 0xff},
		color.RGBA{0xbe, 0x1a, 0x24, 0xff},
		color.RGBA{0x30, 0xe6, 0xc6, 0xff},
		color.RGBA{0xb4, 0x1a, 0xe2, 0xff},
		color.RGBA{0x1f, 0xd2, 0x1e, 0xff},
		color.RGBA{0x21, 0x1b, 0xae, 0xff},
		color.RGBA{0xdf, 0xf6, 0x0a, 0xff},
		color.RGBA{0xb8, 0x41, 0x04, 0xff},
		color.RGBA{0x6a, 0x33, 0x04, 0xff},
		color.RGBA{0xfe, 0x4a, 0x57, 0xff},
		color.RGBA{0x42, 0x45, 0x40, 0xff},
		color.RGBA{0x70, 0x74, 0x6f, 0xff},
		color.RGBA{0x59, 0xfe, 0x59, 0xff},
		color.RGBA{0x5f, 0x53, 0xfe, 0xff},
		color.RGBA{0xa4, 0xa7, 0xa2, 0xff},
	})

var Perplex2 = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xcd, 0x31, 0x00, 0xff},
		color.RGBA{0x81, 0xff, 0xd8, 0xff},
		color.RGBA{0xeb, 0x4c, 0xe1, 0xff},
		color.RGBA{0x69, 0xf7, 0x00, 0xff},
		color.RGBA{0x2a, 0x19, 0xdd, 0xff},
		color.RGBA{0xff, 0xff, 0x5b, 0xff},
		color.RGBA{0xe4, 0x6a, 0x00, 0xff},
		color.RGBA{0x75, 0x48, 0x2a, 0xff},
		color.RGBA{0xff, 0xa1, 0x96, 0xff},
		color.RGBA{0x70, 0x70, 0x70, 0xff},
		color.RGBA{0xb0, 0xb0, 0xb0, 0xff},
		color.RGBA{0xb9, 0xff, 0xb9, 0xff},
		color.RGBA{0x69, 0xc4, 0xff, 0xff},
		color.RGBA{0xeb, 0xeb, 0xeb, 0xff},
	})

var Colodore = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xdb, 0x3a, 0x45, 0xff},
		color.RGBA{0x6c, 0xff, 0xff, 0xff},
		color.RGBA{0xe2, 0x3b, 0xf3, 0xff},
		color.RGBA{0x50, 0xf8, 0x3c, 0xff},
		color.RGBA{0x3f, 0x3a, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x3c, 0xff},
		color.RGBA{0xe2, 0x69, 0x09, 0xff},
		color.RGBA{0x93, 0x55, 0x00, 0xff},
		color.RGBA{0xff, 0x80, 0x8a, 0xff},
		color.RGBA{0x6f, 0x6f, 0x6f, 0xff},
		color.RGBA{0xa6, 0xa6, 0xa6, 0xff},
		color.RGBA{0xb1, 0xff, 0x9f, 0xff},
		color.RGBA{0x91, 0x8b, 0xff, 0xff},
		color.RGBA{0xe1, 0xe1, 0xe1, 0xff},
	})

var Vice371Internal = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0xaf, 0x3c, 0x58, 0xff},
		color.RGBA{0x7e, 0xf3, 0xd6, 0xff},
		color.RGBA{0xaa, 0x40, 0xf5, 0xff},
		color.RGBA{0x62, 0xd5, 0x32, 0xff},
		color.RGBA{0x2c, 0x3d, 0xec, 0xff},
		color.RGBA{0xff, 0xff, 0x46, 0xff},
		color.RGBA{0xb7, 0x63, 0x1e, 0xff},
		color.RGBA{0x77, 0x53, 0x00, 0xff},
		color.RGBA{0xee, 0x7b, 0x95, 0xff},
		color.RGBA{0x62, 0x62, 0x62, 0xff},
		color.RGBA{0x94, 0x94, 0x94, 0xff},
		color.RGBA{0xb7, 0xff, 0x86, 0xff},
		color.RGBA{0x73, 0x85, 0xff, 0xff},
		color.RGBA{0xcd, 0xcd, 0xcd, 0xff},
	})

var Lenyn = createPalette(
	color.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x9d, 0x4b, 0x32, 0xff},
		color.RGBA{0x82, 0xcd, 0xdf, 0xff},
		color.RGBA{0x9d, 0x4a, 0xd7, 0xff},
		color.RGBA{0x72, 0xbe, 0x28, 0xff},
		color.RGBA{0x4a, 0x32, 0xd4, 0xff},
		color.RGBA{0xdd, 0xee, 0x56, 0xff},
		color.RGBA{0xa3, 0x6f, 0x05, 0xff},
		color.RGBA{0x6d, 0x5a, 0x00, 0xff},
		color.RGBA{0xce, 0x84, 0x6f, 0xff},
		color.RGBA{0x64, 0x64, 0x64, 0xff},
		color.RGBA{0x8f, 0x8f, 0x8f, 0xff},
		color.RGBA{0xb6, 0xfb, 0x78, 0xff},
		color.RGBA{0x8a, 0x76, 0xff, 0xff},
		color.RGBA{0xba, 0xba, 0xba, 0xff},
	})

func createPalette(palette color.Palette) Palette {
	result := make(Palette)
	for i, c := range palette {
		result[i] = toColorful(ToRGBA(c))
	}
	return result
}

