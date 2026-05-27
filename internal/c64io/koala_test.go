package c64io

import (
	"testing"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestKoalaBinary(t *testing.T) {
	// 1. Setup a valid Koala-sized image (160x200 multi-color)
	width, height := 160, 200
	spec := indexedimage.Retrospec{
		Name:         indexedimage.KoalaType,
		BitsPerPixel: 2,
	}

	img := &indexedimage.IndexedImage{
		Width:   width,
		Height:  height,
		Spec:    spec,
		Pixels:  make([]pixels.Pixel, width*height),
		Regions: make(map[int][]indexedimage.Region),
	}

	// Mock pixels
	for i := range img.Pixels {
		img.Pixels[i].BitPattern = 0b00 // Default to background
	}

	// Layer 0: Background color (bit pattern 00)
	img.Regions[0] = []indexedimage.Region{
		{
			BitpatternToColor: map[int]int{0b00: 5}, // Background color 5
		},
	}

	// Layer 1: Character regions (4x8 pixels each)
	// 160/4 = 40 cols, 200/8 = 25 rows => 1000 regions
	numRegions := 1000
	img.Regions[1] = make([]indexedimage.Region, numRegions)
	for i := range numRegions {
		// Set specific colors for testing packing:
		// 01 -> 1 (0x1), 10 -> 2 (0x2) => Packed Screen RAM byte should be 0x12
		// 11 -> 3 (0x3) => Color RAM byte should be 0x03
		img.Regions[1][i] = indexedimage.Region{
			BitpatternToColor: map[int]int{
				0b01: 1,
				0b10: 2,
				0b11: 3,
			},
		}
	}

	binFile := koalaBinary(img)
	data := binFile.getBytes()

	// Task 2.1: Verify file size
	if len(data) != 10003 {
		t.Fatalf("Expected 10003 bytes, got %d", len(data))
	}

	// Task 2.2: Verify memory layout order
	// 1. Load Address (2 bytes): 0x6000 -> 00 60
	if data[0] != 0x00 || data[1] != 0x60 {
		t.Errorf("Expected load address 0x6000 (00 60), got %02x %02x", data[0], data[1])
	}

	// 2. Bitmap (8000 bytes)
	// Starts at index 2, ends at 8002

	// 3. Screen RAM (1000 bytes)
	// Starts at index 8002, ends at 9002
	screenRamStart := 8002
	// Task 2.3: Verify correct color packing in Screen RAM (0x12)
	for i := range 1000 {
		if data[screenRamStart+i] != 0x12 {
			t.Errorf("Expected Screen RAM byte 0x12 at index %d, got %02x", screenRamStart+i, data[screenRamStart+i])
			break
		}
	}

	// 4. Color RAM (1000 bytes)
	// Starts at index 9002, ends at 10002
	colorRamStart := 9002
	for i := range 1000 {
		if data[colorRamStart+i] != 0x03 {
			t.Errorf("Expected Color RAM byte 0x03 at index %d, got %02x", colorRamStart+i, data[colorRamStart+i])
			break
		}
	}

	// 5. Background Color (1 byte)
	// Index 10002
	if data[10002] != 0x05 {
		t.Errorf("Expected background color 0x05 at end of file, got %02x", data[10002])
	}
}

func TestKoalaBitmapPacking(t *testing.T) {
	width, height := 160, 200
	spec := indexedimage.Retrospec{
		Name:         indexedimage.KoalaType,
		BitsPerPixel: 2,
	}

	img := &indexedimage.IndexedImage{
		Width:   width,
		Height:  height,
		Spec:    spec,
		Pixels:  make([]pixels.Pixel, width*height),
		Regions: make(map[int][]indexedimage.Region),
	}

	// Set pixel patterns: 00 01 10 11 -> 00011011 in binary = 0x1B
	for i := 0; i < width*height; i += 4 {
		img.Pixels[i].BitPattern = 0b00
		img.Pixels[i+1].BitPattern = 0b01
		img.Pixels[i+2].BitPattern = 0b10
		img.Pixels[i+3].BitPattern = 0b11
	}

	// Mandatory background for koalaBinary
	img.Regions[0] = []indexedimage.Region{{BitpatternToColor: map[int]int{0b00: 0}}}
	img.Regions[1] = make([]indexedimage.Region, 1000)
	for i := range 1000 {
		img.Regions[1][i] = indexedimage.Region{BitpatternToColor: map[int]int{0b01: 0, 0b10: 0, 0b11: 0}}
	}

	binFile := koalaBinary(img)
	data := binFile.getBytes()

	// Bitmap data starts at index 2.
	// Check the first few bytes.
	// Since Koala uses reOrder(CharsByteOrder), the first 8 bytes correspond to the first 4x8 cell.
	// Our pixels were set 00 01 10 11 repeatedly.
	// One byte is 4 pixels. So every byte in the cell should be 0x1B.
	for i := 0; i < 8; i++ {
		if data[2+i] != 0x1B {
			t.Errorf("Expected bitmap byte 0x1b at index %d, got %02x", 2+i, data[2+i])
		}
	}
}
