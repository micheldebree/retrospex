package indexedimage

import (
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestDistance(t *testing.T) {
	color1 := colorful.Color{R: 0.5, G: 0.5, B: 0.5}
	color2 := colorful.Color{R: 0.7, G: 0.7, B: 0.7}

	expectedDistance := color1.DistanceLinearRGB(color2)
	result := distance(color1, color2)

	if result != expectedDistance {
		t.Errorf("Expected distance %f, got %f", expectedDistance, result)
	}
}

func TestDistances(t *testing.T) {
	color := colorful.Color{R: 0.5, G: 0.5, B: 0.5}
	palette := pixels.Palette{
		0: colorful.Color{R: 0.1, G: 0.1, B: 0.1},
		1: colorful.Color{R: 0.9, G: 0.9, B: 0.9},
	}

	result := distances(color, palette)

	if len(result) != len(palette) {
		t.Errorf("Expected %d distances, got %d", len(palette), len(result))
	}
}

func TestBestPixelIndex(t *testing.T) {
	distances := PaletteDistance{
		0: 0.1,
		1: 0.2,
		2: 0.3,
	}

	index, qerror := bestPixelIndex(distances)

	if index != 0 || qerror != 0.1 {
		t.Errorf("Expected index %d and qerror %f, got index %d and qerror %f", 0, 0.1, index, qerror)
	}
}

func TestQuantizeToIndex(t *testing.T) {
	color := colorful.Color{R: 0.5, G: 0.5, B: 0.5}
	palette := pixels.Palette{
		0: colorful.Color{R: 0.1, G: 0.1, B: 0.1},
		1: colorful.Color{R: 0.9, G: 0.9, B: 0.9},
	}

	index, qerror := QuantizeToIndex(color, palette)

	if index != 1 || qerror > 0.5 {
		t.Errorf("Expected index %d and qerror <= %f, got index %d and qerror %f", 1, 0.5, index, qerror)
	}
}

func TestQuantizePixel(t *testing.T) {
	pixel := &pixels.Pixel{Color: colorful.Color{R: 0.5, G: 0.5, B: 0.5}}
	palette := pixels.Palette{
		0: colorful.Color{R: 0.1, G: 0.1, B: 0.1},
		1: colorful.Color{R: 0.9, G: 0.9, B: 0.9},
	}

	quantizePixel(pixel, palette)

	if pixel.PaletteIndex != 1 || pixel.QuantizationError > 0.5 {
		t.Errorf("Expected PaletteIndex %d and QuantizationError <= %f, got PaletteIndex %d and QuantizationError %f", 1, 0.5, pixel.PaletteIndex, pixel.QuantizationError)
	}
}

func TestReducePalette(t *testing.T) {
	img := IndexedImage{
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 0.1, G: 0.1, B: 0.1}},
			{Color: colorful.Color{R: 0.9, G: 0.9, B: 0.9}},
			{Color: colorful.Color{R: 0.5, G: 0.5, B: 0.5}},
		},
		palette: pixels.Palette{
			0: colorful.Color{R: 0.1, G: 0.1, B: 0.1},
			1: colorful.Color{R: 0.9, G: 0.9, B: 0.9},
			2: colorful.Color{R: 0.5, G: 0.5, B: 0.5},
		},
	}
	layer := Layer{
		bitpatterns: []int{0, 1},
	}

	result := reducePalette(img, layer)

	if len(result.palette) != 2 || len(result.bitpatterns) != 2 {
		t.Errorf("Expected palette and bitpatterns of length %d, got palette of length %d and bitpatterns of length %d", 2, len(result.palette), len(result.bitpatterns))
	}
}

func TestQuantize(t *testing.T) {
	img := IndexedImage{
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 0.1, G: 0.1, B: 0.1}},
			{Color: colorful.Color{R: 0.9, G: 0.9, B: 0.9}},
			{Color: colorful.Color{R: 0.5, G: 0.5, B: 0.5}},
		},
		palette: pixels.Palette{
			0: colorful.Color{R: 0.1, G: 0.1, B: 0.1},
			1: colorful.Color{R: 0.9, G: 0.9, B: 0.9},
			2: colorful.Color{R: 0.5, G: 0.5, B: 0.5},
		},
		spec: Retrospec{
			layers: []Layer{
				{bitpatterns: []int{0, 1}},
			},
		},
	}

	result := Quantize(img)

	if len(result.pixels) != len(img.pixels) {
		t.Errorf("Expected %d pixels, got %d", len(img.pixels), len(result.pixels))
	}
}
