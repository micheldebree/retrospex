package indexedimage

type Layer struct {
	cellWidth, cellHeight int
	bitpatterns           []int
	// TODO: this can be determined while processing
	isLast bool // the last layer should quantize all remaining pixels
}

type Retrospec struct {
	layers []Layer
}
