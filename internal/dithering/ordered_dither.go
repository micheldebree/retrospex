package dithering

import (
	"image/color"

	"github.com/micheldebree/retrospex/internal/indexedimage"
)

type IntMatrix [][]int

var DitherMatrices = map[string]IntMatrix{
	"none": {{0}},
	"bayer2x2": {
		{0, 2},
		{3, 1},
	},
	"bayer4x4": {
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	},
	"bayer8x8": {
		{1, 49, 13, 61, 4, 52, 16, 64},
		{33, 17, 45, 29, 36, 20, 48, 31},
		{9, 57, 5, 53, 12, 60, 8, 56},
		{41, 25, 37, 21, 44, 28, 40, 24},
		{3, 51, 15, 63, 2, 50, 14, 62},
		{35, 19, 47, 31, 34, 18, 46, 30},
		{11, 59, 7, 55, 10, 58, 6, 54},
		{43, 27, 39, 23, 42, 26, 38, 22},
	},
}

// OrderedDither Apply ordered dithering by offsetting color channels with
// a matrix. N.B. the image itself is adjusted
func OrderedDither(img *indexedimage.IndexedImage, matrix IntMatrix, depth int) {

	// no need to waste time if the effect is none
	if depth == 0 || len(matrix) <= 1 {
		return
	}

	normalizedMatrix := normalize(matrix, depth)
	matrixH := len(matrix)
	matrixW := len(matrix[0])

	for i := range (*img).Pixels {
		// change pixels in place
		pixel := &(*img).Pixels[i]
		matrixX := pixel.X % matrixW
		matrixY := pixel.Y % matrixH
		matrixV := normalizedMatrix[matrixY][matrixX]

		r, g, b, _ := pixel.Color.RGBA()

		rOffset := toByteWithOverflow(r) + matrixV
		gOffset := toByteWithOverflow(g) + matrixV
		bOffset := toByteWithOverflow(b) + matrixV

		pixel.Color = color.RGBA{capToByte(rOffset), capToByte(gOffset), capToByte(bOffset), 0xff}
	}
}

// Scale back to the range of a byte, but leave room for underflow and overflow
// by using int
func toByteWithOverflow(colorChannel uint32) int {
	return int(colorChannel >> 8)
}

// Convert to byte, capping underflow and overflow resulting from adding the
// matrix
func capToByte(value int) uint8 {
	if value > 255 {
		return 255
	}
	if value < 0 {
		return 0
	}
	return uint8(value)
}

func normalize(matrix IntMatrix, depth int) IntMatrix {

	height := len(matrix)
	width := len(matrix[0])

	// initialize float matrix
	result := make(IntMatrix, height)
	for i := range result {
		result[i] = make([]int, width)
	}

	factor := float64(depth) / float64(width*height)
	for rowIndex, row := range matrix {
		for colIndex, col := range row {
			result[rowIndex][colIndex] = int(factor*float64(col) - 0.5)
		}
	}
	return result
}
