package c64io

import (
	"testing"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestGetBitmapData(t *testing.T) {
	tests := []struct {
		name     string
		image    *indexedimage.IndexedImage
		expected []byte
	}{
		{
			name: "Single pixel",
			image: &indexedimage.IndexedImage{
				Width:  1,
				Height: 1,
				Spec: indexedimage.Retrospec{
					BitPatternSize: 1,
				},
				pixels: []pixels.Pixel{
					{
						X:         0,
						Y:         0,
						Color:     pixels.ToColorful(color.RGBA{0xff, 0x00, 0x00, 0xff}),
						BitPattern: 1,
					},
				},
			},
			expected: []byte{0x80}, // 1 bit pattern (1) in the most significant position
		},
		{
			name: "Multiple pixels",
			image: &indexedimage.IndexedImage{
				Width:  2,
				Height: 2,
				Spec: indexedimage.Retrospec{
					BitPatternSize: 1,
				},
				pixels: []pixels.Pixel{
					{X: 0, Y: 0, Color: pixels.ToColorful(color.RGBA{0xff, 0x00, 0x00, 0xff}), BitPattern: 1},
					{X: 1, Y: 0, Color: pixels.ToColorful(color.RGBA{0x00, 0xff, 0x00, 0xff}), BitPattern: 2},
					{X: 0, Y: 1, Color: pixels.ToColorful(color.RGBA{0x00, 0x00, 0xff, 0xff}), BitPattern: 3},
					{X: 1, Y: 1, Color: pixels.ToColorful(color.RGBA{0xff, 0xff, 0x00, 0xff}), BitPattern: 4},
				},
			},
			expected: []byte{0xc0, 0xf0}, // 0b11000000, 0b11110000
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
					t.Errorf("at index %d: expected %v, got %v", i, tt.expected[i], v)
				}
			}
		})
	}

	t.Run("Pixel without bit pattern", func(t *testing.T) {
		image := &indexedimage.IndexedImage{
			Width:  1,
			Height: 1,
			Spec: indexedimage.Retrospec{
				BitPatternSize: 1,
			},
			pixels: []pixels.Pixel{
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
			result := reOrderToC64BitmapOrder(tt.input, tt.bytesPerInputRow)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i, tt.expected[i], v)
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
		reOrderToC64BitmapOrder(input, bytesPerInputRow)
	})
}
