package main

type Layer struct {
	cellWidth, cellHeight int
	bitpatterns           []int
	// TODO: this can be determined while processing
	isLast bool // the last layer should quantize all remaining pixels
}

type Retrospec struct {
	layers []Layer
}

var Hires2Colors = Retrospec{
	[]Layer{
		{320, 200, []int{0, 1}, true}, // 0400,x (lower), 0400,x (higher)
	},
}
