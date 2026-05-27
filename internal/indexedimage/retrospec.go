package indexedimage

import "fmt"

type Layer struct {
	CellWidth, CellHeight int
	Bitpatterns           []int
	// TODO: this can be determined while processing
	IsLast bool // the last layer should quantize all remaining pixels
}

type LayerTemplate struct {
	UseImageDimensions    bool // if true, CellWidth/Height will be set from image dimensions
	CellWidth, CellHeight int  // used when UseImageDimensions is false
	Bitpatterns           []int
	IsLast                bool // the last layer should quantize all remaining pixels
}

type RetrospecName string

type Retrospec struct {
	Name         RetrospecName
	Layers       []Layer
	BitsPerPixel int
}

type RetrospecTemplate struct {
	Name         RetrospecName
	Layers       []LayerTemplate
	BitsPerPixel int
}

func (spec Retrospec) PixelsPerByte() int {
	return 8 / spec.BitsPerPixel
}

func (spec Retrospec) Print() {
	fmt.Printf("Mode: %s\nBits per pixel: %d\n", spec.Name, spec.BitsPerPixel)
}
