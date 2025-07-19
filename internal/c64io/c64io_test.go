package c64io

import (
	"testing"

	"image/color"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestGetBitmapData(t *testing.T) {

	irrelevantColor := pixels.ToColorful(color.RGBA{0, 0, 0, 0})

	tests := []struct {
		name     string
		image    *indexedimage.IndexedImage
		expected []byte
	}{
		{
			name: "Multicolor",
			image: &indexedimage.IndexedImage{
				Width:  8,
				Height: 2,
				Spec: indexedimage.Retrospec{
					BitPatternSize: 2,
				},
				Pixels: []pixels.Pixel{
					{X: 0, Y: 0, Color: irrelevantColor, BitPattern: 0b01},
					{X: 1, Y: 0, Color: irrelevantColor, BitPattern: 0b10},
					{X: 2, Y: 0, Color: irrelevantColor, BitPattern: 0b11},
					{X: 3, Y: 0, Color: irrelevantColor, BitPattern: 0b10},
					{X: 4, Y: 0, Color: irrelevantColor, BitPattern: 0b00},
					{X: 5, Y: 0, Color: irrelevantColor, BitPattern: 0b11},
					{X: 6, Y: 0, Color: irrelevantColor, BitPattern: 0b11},
					{X: 7, Y: 0, Color: irrelevantColor, BitPattern: 0b01},
					{X: 0, Y: 1, Color: irrelevantColor, BitPattern: 0b00},
					{X: 1, Y: 1, Color: irrelevantColor, BitPattern: 0b11},
					{X: 2, Y: 1, Color: irrelevantColor, BitPattern: 0b11},
					{X: 3, Y: 1, Color: irrelevantColor, BitPattern: 0b01},
					{X: 4, Y: 1, Color: irrelevantColor, BitPattern: 0b10},
					{X: 5, Y: 1, Color: irrelevantColor, BitPattern: 0b10},
					{X: 6, Y: 1, Color: irrelevantColor, BitPattern: 0b01},
					{X: 7, Y: 1, Color: irrelevantColor, BitPattern: 0b00},
				},
			},
			expected: []byte{
				0b01101110, 0b00111101,
				0b00111101, 0b10100100,
			},
		},

		{
			name: "Singlecolor",
			image: &indexedimage.IndexedImage{
				Width:  16,
				Height: 2,
				Spec: indexedimage.Retrospec{
					BitPatternSize: 1,
				},
				Pixels: []pixels.Pixel{
					{X: 0, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 1, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 2, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 3, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 4, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 5, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 6, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 7, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 8, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 9, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 10, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 11, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 12, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 13, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 14, Y: 0, Color: irrelevantColor, BitPattern: 0b0},
					{X: 15, Y: 0, Color: irrelevantColor, BitPattern: 0b1},
					{X: 0, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 1, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 2, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 3, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 4, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 5, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 6, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 7, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 8, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 9, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 10, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 11, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 12, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 13, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
					{X: 14, Y: 1, Color: irrelevantColor, BitPattern: 0b1},
					{X: 15, Y: 1, Color: irrelevantColor, BitPattern: 0b0},
				},
			},
			expected: []byte{
				0b01010101, 0b01010101,
				0b10101010, 0b10101010,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBitmapData(tt.image)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %08b, got %08b", i, tt.expected[i], v)
				}
			}
		})
	}

	t.Run("Pixel without bit pattern", func(t *testing.T) {
		image := &indexedimage.IndexedImage{
			Width:  8,
			Height: 8,
			Spec: indexedimage.Retrospec{
				BitPatternSize: 1,
			},
			Pixels: []pixels.Pixel{
				{X: 0, Y: 0, Color: pixels.ToColorful(color.RGBA{0xff, 0x00, 0x00, 0xff}), BitPattern: -1}, // No bit pattern
			},
		}
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		getBitmapData(image)
	})
}

func TestReOrderToC64BitmapOrder(t *testing.T) {
	tests := []struct {
		name             string
		input            []byte
		bytesPerInputRow int
		expected         []byte
	}{
		{
			name: "Single row",
			input: []byte{
				0x01, 0x02,
				0x03, 0x04,
				0x05, 0x06,
				0x07, 0x08,
				0x09, 0x0a,
				0x0b, 0x0c,
				0x0d, 0x0e,
				0x0f, 0x10,
			},
			bytesPerInputRow: 2,
			expected: []byte{
				0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f,
				0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e, 0x10,
			},
		},
		{
			name: "Multiple rows",
			input: []byte{
				0x01, 0x02,
				0x03, 0x04,
				0x05, 0x06,
				0x07, 0x08,
				0x09, 0x0a,
				0x0b, 0x0c,
				0x0d, 0x0e,
				0x0f, 0x10,
				0x11, 0x12,
				0x13, 0x14,
				0x15, 0x16,
				0x17, 0x18,
				0x19, 0x1a,
				0x1b, 0x1c,
				0x1d, 0x1e,
				0x1f, 0x20,
			},
			bytesPerInputRow: 2,
			expected: []byte{
				0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f,
				0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e, 0x10,
				0x11, 0x13, 0x15, 0x17, 0x19, 0x1b, 0x1d, 0x1f,
				0x12, 0x14, 0x16, 0x18, 0x1a, 0x1c, 0x1e, 0x20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reOrderToVicBitmapOrder(tt.input, tt.bytesPerInputRow)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %08b, got %08b", i, tt.expected[i], v)
				}
			}
		})
	}

	t.Run("Input height not multiple of 8", func(t *testing.T) {
		input := []byte{
			0x01, 0x02, 0x03, 0x04,
			0x11, 0x12, 0x13, 0x14,
		}
		bytesPerInputRow := 8
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		reOrderToVicBitmapOrder(input, bytesPerInputRow)
	})
}
