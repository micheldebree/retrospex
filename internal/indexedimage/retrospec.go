package indexedimage

type Layer struct {
	CellWidth, CellHeight int
	Bitpatterns           []int
	// TODO: this can be determined while processing
	IsLast bool // the last layer should quantize all remaining pixels
}

type RetrospecName string

type Retrospec struct {
	Name           RetrospecName
	Layers         []Layer
	BitPatternSize int
}

func (spec Retrospec) PixelsPerByte() int {
	return 8 / spec.BitPatternSize
}
