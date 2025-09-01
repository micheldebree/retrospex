package indexedimage

import "fmt"

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

func (retrospec Retrospec) Print() {
	fmt.Printf("Mode: %s\nBits per pixel: %d\n", retrospec.Name, retrospec.BitPatternSize)
}
