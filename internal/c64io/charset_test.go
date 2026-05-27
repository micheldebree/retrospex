package c64io

import (
	"testing"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestCharsetBinary(t *testing.T) {
	// Setup a 64x64 image (8x8 = 64 characters)
	width, height := 64, 64
	spec := indexedimage.Retrospec{
		Name:         indexedimage.SCCharsetType,
		BitsPerPixel: 1,
	}

	img := &indexedimage.IndexedImage{
		Width:  width,
		Height: height,
		Spec:   spec,
		Pixels: make([]pixels.Pixel, width*height),
	}

	// Set alternating bit patterns
	for i := range img.Pixels {
		img.Pixels[i].BitPattern = i % 2
	}

	binFile := charsetBinary(img)
	data := binFile.getBytes()

	// 64x64 pixels = 4096 pixels. 1 bit per pixel = 4096 bits = 512 bytes.
	// 512 bytes / 8 bytes per char = 64 characters.
	if len(data) != 512 {
		t.Fatalf("Expected 512 bytes, got %d", len(data))
	}
}

func TestMCCharsetBinary(t *testing.T) {
	// Setup a 8x8 image. 
	// In MCCharset, width 8 pixels (squashed) means 2 characters of width 4.
	// Height 8 pixels means 1 row of characters.
	// Total 2 characters.

	width, height := 8, 8
	palette := pixels.C64Palettes["colodore"]
	spec := indexedimage.Retrospec{
		Name: indexedimage.MCCharsetType,
		Layers: []indexedimage.Layer{
			{CellWidth: width, CellHeight: height, Bitpatterns: []int{0b00, 0b01, 0b10}, IsLast: false},
			{CellWidth: 4, CellHeight: 8, Bitpatterns: []int{0b11}, IsLast: true},
		},
		BitsPerPixel: 2,
	}

	img := &indexedimage.IndexedImage{
		Width:   width,
		Height:  height,
		Spec:    spec,
		Palette: palette,
		Pixels:  make([]pixels.Pixel, width*height),
		Regions: make(map[int][]indexedimage.Region),
	}

	// Layer 0 region (global colors)
	img.Regions[0] = []indexedimage.Region{
		{
			BitpatternToColor: map[int]int{0b00: 0, 0b01: 1, 0b10: 2},
		},
	}

	// Layer 1 regions (2 regions)
	img.Regions[1] = []indexedimage.Region{
		{BitpatternToColor: map[int]int{0b11: 3}},
		{BitpatternToColor: map[int]int{0b11: 4}},
	}

	binFile := charsetBinary(img)
	data := binFile.getBytes()

	// 2 characters.
	// Bitmap: 2 * 8 = 16 bytes.
	// D800: 2 * 1 = 2 bytes.
	// Global: 3 bytes.
	// Total: 21 bytes.
	expectedSize := (9 * 2) + 3
	if len(data) != expectedSize {
		t.Fatalf("Expected %d bytes, got %d", expectedSize, len(data))
	}

	// Verify global colors at the end
	globalStart := 16 + 2
	if data[globalStart] != 0 || data[globalStart+1] != 1 || data[globalStart+2] != 2 {
		t.Errorf("Expected global colors 0, 1, 2, got %v", data[globalStart:globalStart+3])
	}

	// Verify D800 colors in the middle
	d800Start := 16
	if data[d800Start] != 3 || data[d800Start+1] != 4 {
		t.Errorf("Expected D800 colors 3, 4, got %v", data[d800Start:d800Start+2])
	}
}
