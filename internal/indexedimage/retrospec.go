package indexedimage

type Layer struct {
	CellWidth, CellHeight int
	Bitpatterns           []int
	// TODO: this can be determined while processing
	IsLast bool // the last layer should quantize all remaining pixels
}

type Retrospec struct {
	Layers         []Layer
	BitPatternSize int
}
