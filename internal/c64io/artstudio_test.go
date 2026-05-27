package c64io

import (
	"testing"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestArtstudioBinary(t *testing.T) {
	// Setup a 320x200 hires image
	width, height := 320, 200
	spec := indexedimage.Retrospec{
		Name:         indexedimage.HiresType,
		BitsPerPixel: 1,
	}

	img := &indexedimage.IndexedImage{
		Width:  width,
		Height: height,
		Spec:   spec,
		Pixels: make([]pixels.Pixel, width*height),
		Regions: map[int][]indexedimage.Region{
			0: make([]indexedimage.Region, 1000), // 40x25 = 1000 regions
		},
	}

	// Mock pixels
	for i := range img.Pixels {
		img.Pixels[i].BitPattern = 0
	}

	for i := range 1000 {
		// Bit 1: color 1, Bit 0: color 0 => Packed byte should be 0x10
		img.Regions[0][i] = indexedimage.Region{
			BitpatternToColor: map[int]int{
				1: 1,
				0: 0,
			},
		}
	}

	binFile := artstudioBinary(img)
	data := binFile.getBytes()

	// 2 (address) + 8000 (bitmap) + 1000 (screen) + 7 (padding) = 9009 bytes
	if len(data) != 9009 {
		t.Fatalf("Expected 9009 bytes, got %d", len(data))
	}

	// Verify Load Address: 0x2000 -> 00 20
	if data[0] != 0x00 || data[1] != 0x20 {
		t.Errorf("Expected load address 0x2000 (00 20), got %02x %02x", data[0], data[1])
	}

	// Screen RAM starts after address (2) + bitmap (8000) = 8002
	screenRamStart := 8002
	for i := range 1000 {
		if data[screenRamStart+i] != 0x10 {
			t.Errorf("Expected Screen RAM byte 0x10 at index %d, got %02x", screenRamStart+i, data[screenRamStart+i])
			break
		}
	}

	// Padding (7 bytes) at the end
	for i := 0; i < 7; i++ {
		if data[9002+i] != 0 {
			t.Errorf("Expected padding byte 0 at index %d, got %02x", 9002+i, data[9002+i])
		}
	}
}
