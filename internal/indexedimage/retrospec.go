package indexedimage

type Layer struct {
	CellWidth, CellHeight int
	Bitpatterns           []int
	// TODO: this can be determined while processing
	IsLast bool // the last layer should quantize all remaining pixels
}

type RetrospecName string

type ByteOrder int

const (
	RegularByteOrder ByteOrder = iota // regular, row-first byte order
	VicByteOrder                      // quirky vic byte order
)

type Retrospec struct {
	Name            RetrospecName
	Layers          []Layer
	BitPatternSize  int
	BitmapByteOrder ByteOrder
}

func (spec Retrospec) PixelsPerByte() int {
	return 8 / spec.BitPatternSize
}
