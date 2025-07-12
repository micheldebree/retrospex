package indexedimage

import (
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestDistance(t *testing.T) {
	color1 := colorful.Color{R: 1, G: 0, B: 0}
	color2 := colorful.Color{R: 0, G: 1, B: 0}
	expected := color1.DistanceLinearRGB(color2)
	actual := distance(color1, color2)
	if actual != expected {
		t.Errorf("Expected distance %f but got %f", expected, actual)
	}
}

func TestDistances(t *testing.T) {
	color := colorful.Color{R: 0.5, G: 0.5, B: 0.5}
	palette := pixels.Palette{
		0: colorful.Color{R: 1, G: 0, B: 0},
		1: colorful.Color{R: 0, G: 1, B: 0},
		2: colorful.Color{R: 0, G: 0, B: 1},
	}
	expected := PaletteDistance{
		0: color.DistanceLinearRGB(colorful.Color{R: 1, G: 0, B: 0}),
		1: color.DistanceLinearRGB(colorful.Color{R: 0, G: 1, B: 0}),
		2: color.DistanceLinearRGB(colorful.Color{R: 0, G: 0, B: 1}),
	}
	actual := distances(color, palette)
	for index, expectedDist := range expected {
		if actual[index] != expectedDist {
			t.Errorf("Expected distance %f for index %d but got %f", expectedDist, index, actual[index])
		}
	}
}

func TestBestPixelIndex(t *testing.T) {
	distances := PaletteDistance{
		0: 1.0,
		1: 0.5,
		2: 0.7,
	}
	expectedIndex := 1
	expectedError := 0.5
	index, error := bestPixelIndex(distances)
	if index != expectedIndex || error != expectedError {
		t.Errorf("Expected best pixel index %d with error %f but got %d with error %f", expectedIndex, expectedError, index, error)
	}
}

func TestQuantizePixel(t *testing.T) {
	pixel := &pixels.Pixel{Color: colorful.Color{R: 0.5, G: 0.5, B: 0.5}}
	palette := pixels.Palette{
		0: colorful.Color{R: 1, G: 0, B: 0},
		1: colorful.Color{R: 0, G: 1, B: 0},
		2: colorful.Color{R: 0, G: 0, B: 1},
	}
	expectedIndex := 1
	quantizePixel(pixel, palette)
	if pixel.PaletteIndex != expectedIndex {
		t.Errorf("Expected quantized index %d but got %d", expectedIndex, pixel.PaletteIndex)
	}
}

func TestLayerRegion_addMapping(t *testing.T) {
	region := LayerRegion{
		bitpatternToColor: make(map[int]int),
		colorToBitpattern: make(map[int]int),
	}
	expectedIndex := 0
	region.addMapping(1, expectedIndex)
	if _, exists := region.bitpatternToColor[1]; !exists {
		t.Error("Expected bitpattern to be added to map")
	}
	if index, _ := region.colorToBitpattern[expectedIndex]; index != 1 {
		t.Error("Expected color to be mapped to bitpattern in reverse map")
	}
}

func TestLayerRegion_coordsToIndex(t *testing.T) {
	region := LayerRegion{img: &IndexedImage{width: 3, height: 3}}
	expectedIndex := 4
	x, y := 1, 1
	index := region.coordsToIndex(x, y)
	if index != expectedIndex {
		t.Errorf("Expected index %d for coords (%d,%d) but got %d", expectedIndex, x, y, index)
	}
}

func TestLayerRegion_getUnmappedBitPattern(t *testing.T) {
	region := LayerRegion{
		bitpatternToColor: map[int]int{1: -1, 2: -1},
		colorToBitpattern: map[int]int{},
	}
	expectedPattern := 1
	pattern := region.getUnmappedBitPattern()
	if pattern != expectedPattern {
		t.Errorf("Expected unmapped bitpattern %d but got %d", expectedPattern, pattern)
	}
}

func TestLayerRegion_getPalette(t *testing.T) {
	region := LayerRegion{
		img: &IndexedImage{
			pixels: []pixels.Pixel{
				{Color: colorful.Color{R: 1, G: 0, B: 0}, PaletteIndex: 0},
				{Color: colorful.Color{R: 0, G: 1, B: 0}, PaletteIndex: 1},
			},
		},
		bitpatternToColor: map[int]int{1: 0, 2: -1},
	}
	expectedPalette := pixels.Palette{0: colorful.Color{R: 1, G: 0, B: 0}}
	palette := region.getPalette()
	for index, color := range expectedPalette {
		if palette[index] != color {
			t.Errorf("Expected color %v at index %d in palette but got %v", color, index, palette[index])
		}
	}
}

func TestGetLayerRegions(t *testing.T) {
	img := IndexedImage{
		width:  3,
		height: 3,
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 1, G: 0, B: 0}},
			{Color: colorful.Color{R: 0, G: 1, B: 0}},
		},
	}
	layer := Layer{cellWidth: 1, cellHeight: 1}
	expectedRegions := []LayerRegion{
		{img: &img, x: 0, y: 0, width: 1, height: 1, bitpatternToColor: map[int]int{1: -1}, colorToBitpattern: map[int]int{}, isLastLayer: false},
		{img: &img, x: 1, y: 0, width: 1, height: 1, bitpatternToColor: map[int]int{2: -1}, colorToBitpattern: map[int]int{}, isLastLayer: false},
		{img: &img, x: 0, y: 1, width: 1, height: 1, bitpatternToColor: map[int]int{-1, -1}, colorToBitpattern: map[int]int{}},
		{img: &img, x: 1, y: 1, width: 1, height: 1, bitpatternToColor: map[int]int{-1, -1}, colorToBitpattern: map[int]int{}},
	}
	regions := getLayerRegions(img, layer)
	for i, expectedRegion := range expectedRegions {
		if !compareLayerRegions(&expectedRegion, &regions[i]) {
			t.Errorf("Expected region %v but got %v", expectedRegion, regions[i])
		}
	}
}

func compareLayerRegions(r1 *LayerRegion, r2 *LayerRegion) bool {
	if r1.x != r2.x || r1.y != r2.y || r1.width != r2.width || r1.height != r2.height {
		return false
	}
	for pattern, index := range r1.bitpatternToColor {
		if r2.bitpatternToColor[pattern] != index {
			return false
		}
	}
	for index, pattern := range r1.colorToBitpattern {
		if r2.colorToBitpattern[index] != pattern {
			return false
		}
	}
	return true
}

func TestQuantizeLayerRegion(t *testing.T) {
	img := IndexedImage{
		width:  3,
		height: 3,
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 1, G: 0, B: 0}},
			{Color: colorful.Color{R: 0, G: 1, B: 0}},
		},
		palette: pixels.Palette{
			0: colorful.Color{R: 1, G: 0, B: 0},
			1: colorful.Color{R: 0, G: 1, B: 0},
		},
	}
	layer := Layer{cellWidth: 1, cellHeight: 1}
	regions := getLayerRegions(img, layer)
	for i := range regions {
		quantizeLayerRegion(regions[i])
		if !regions[i].img.pixels[i].isQuantized() {
			t.Errorf("Expected pixel %v to be quantized", regions[i].img.pixels[i])
		}
	}
}

func TestQuantize(t *testing.T) {
	img := IndexedImage{
		width:  3,
		height: 3,
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 1, G: 0, B: 0}},
			{Color: colorful.Color{R: 0, G: 1, B: 0}},
		},
		palette: pixels.Palette{
			0: colorful.Color{R: 1, G: 0, B: 0},
			1: colorful.Color{R: 0, G: 1, B: 0},
		},
	}
	expectedImg := IndexedImage{
		width:  3,
		height: 3,
		pixels: []pixels.Pixel{
			{Color: colorful.Color{R: 1, G: 0, B: 0}, PaletteIndex: 0},
			{Color: colorful.Color{R: 0, G: 1, B: 0}, PaletteIndex: 1},
		},
	}
	result := Quantize(img)
	if !compareIndexedImages(&expectedImg, &result) {
		t.Errorf("Expected quantized image %v but got %v", expectedImg, result)
	}
}

func compareIndexedImages(i1 *IndexedImage, i2 *IndexedImage) bool {
	if i1.width != i2.width || i1.height != i2.height {
		return false
	}
	for