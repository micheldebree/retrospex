package c64io

import (
	"testing"
)

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
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			},
			bytesPerInputRow: 8,
			expected: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			},
		},
		{
			name: "Multiple rows",
			input: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			},
			bytesPerInputRow: 8,
			expected: []byte{
				0x01, 0x11, 0x02, 0x12, 0x03, 0x13, 0x04, 0x14,
				0x05, 0x15, 0x06, 0x16, 0x07, 0x17, 0x08, 0x18,
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
