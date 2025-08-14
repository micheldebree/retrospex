package conversion

import (
	"testing"

	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
	"github.com/stretchr/testify/assert"
)

func TestRegion_coordsToIndex(t *testing.T) {
	img := &indexedimage.IndexedImage{
		Width:   10,
		Height:  10,
		Spec:    indexedimage.Retrospec{},
		Palette: pixels.Palette{},
		Pixels:  []pixels.Pixel{},
	}
	region := &Region{
		img:               img,
		x:                 2,
		y:                 3,
		width:             10,
		height:            10,
		bitpatternToColor: make(map[int]int),
		colorToBitpattern: make(map[int]int),
		isLastLayer:       false,
	}

	assert.Equal(t, 32, region.coordsToIndex(0, 0))
	assert.Equal(t, 72, region.coordsToIndex(10, 3))
}

func TestRegion_getUnmappedBitPattern(t *testing.T) {
	img := &indexedimage.IndexedImage{
		Width:   10,
		Height:  10,
		Spec:    indexedimage.Retrospec{},
		Palette: pixels.Palette{},
		Pixels:  []pixels.Pixel{},
	}
	region := &Region{
		img:               img,
		x:                 0,
		y:                 0,
		width:             10,
		height:            10,
		bitpatternToColor: map[int]int{1: -1, 2: 3},
		colorToBitpattern: make(map[int]int),
		isLastLayer:       false,
	}

	found, bitpattern := region.getUnmappedBitPattern()
	assert.True(t, found)
	assert.Equal(t, 1, bitpattern)
}
