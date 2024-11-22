package main

import "image"

func makeKoalaSpec(img *image.Image) Retrospec {

	w, h := getDimensions(img)

	return Retrospec{
		[]Layer{
			{w, h, []int{0x00}, false},            // d021
			{4, 8, []int{0x01, 0x10, 0x11}, true}, // 0400,x (upper nibble), 0400,x (lower nibble), d800,x
		},
	}

}
