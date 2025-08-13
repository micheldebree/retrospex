package conversion

import (
	"image/color"
	"testing"

	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestQuantizePixel(t *testing.T) {
	tests := []struct {
		name     string
		pixel    pixels.Pixel
		palette  pixels.Palette
		expected int
	}{
		{
			name: "Red pixel quantizes to red",
			pixel: pixels.Pixel{
				Color: color.RGBA{0xff, 0x00, 0x00, 0xff},
			},
			palette: pixels.Palette{
				0: color.RGBA{0xff, 0x00, 0x00, 0xff},
				1: color.RGBA{0x00, 0xff, 0x00, 0xff},
			},
			expected: 0,
		},
		{
			name: "Green pixel quantizes to green",
			pixel: pixels.Pixel{
				Color: color.RGBA{0x00, 0xff, 0x00, 0xff},
			},
			palette: pixels.Palette{
				0: color.RGBA{0xff, 0x00, 0x00, 0xff},
				1: color.RGBA{0x00, 0xff, 0x00, 0xff},
			},
			expected: 1,
		},
		{
			name: "Blue pixel quantizes to closest color",
			pixel: pixels.Pixel{
				Color: color.RGBA{0x00, 0x00, 0xff, 0xff},
			},
			palette: pixels.Palette{
				0: color.RGBA{0xff, 0x00, 0x00, 0xff},
				1: color.RGBA{0x00, 0xff, 0x00, 0xff},
				2: color.RGBA{0x00, 0x00, 0xff, 0xff},
			},
			expected: 2,
		},
		{
			name: "Color red and green bit slightly closer to green quantizes to green",
			pixel: pixels.Pixel{
				Color: color.RGBA{0x80, 0x81, 0x00, 0xff},
			},
			palette: pixels.Palette{
				0: color.RGBA{0xff, 0x00, 0x00, 0xff},
				1: color.RGBA{0x00, 0xff, 0x00, 0xff},
			},
			expected: 1,
		},
		{
			name: "Color slightly closer to red than green",
			pixel: pixels.Pixel{
				Color: color.RGBA{0xa0, 0x50, 0x00, 0xff},
			},
			palette: pixels.Palette{
				0: color.RGBA{0xff, 0x00, 0x00, 0xff},
				1: color.RGBA{0x00, 0xff, 0x00, 0xff},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QuantizePixel(&tt.pixel, tt.palette)
			if tt.pixel.PaletteIndex != tt.expected {
				t.Errorf("Expected PaletteIndex to be %d, got %d", tt.expected, tt.pixel.PaletteIndex)
			}
		})
	}
}
