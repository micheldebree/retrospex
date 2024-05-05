package main

import "github.com/lucasb-eyer/go-colorful"

type IntMatrix [][]int
type FloatMatrix [][]float64

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
func OrderedDither(img *IndexedImage, matrix IntMatrix, depth int) {

	normalizedMatrix := normalize(matrix, float64(depth)/255.)
	matrixH := len(matrix)
	matrixW := len(matrix[0])

	for i := range (*img).pixels {
		// change pixels in place
		pixel := &(*img).pixels[i]
		matrixX := pixel.x % matrixW
		matrixY := pixel.y % matrixH
		matrixV := normalizedMatrix[matrixY][matrixX]

		color := colorful.Color{
			R: max(pixel.color.R+matrixV, 0.0),
			G: max(pixel.color.G+matrixV, 0.0),
			B: max(pixel.color.B+matrixV, 0.0),
		}
		pixel.color = color
	}
}

func normalize(matrix IntMatrix, depth float64) FloatMatrix {

	height := len(matrix)
	width := len(matrix[0])

	// initialize float matrix
	result := make(FloatMatrix, height)
	for i := range result {
		result[i] = make([]float64, width)
	}

	factor := 1.0 / float64(width*height)
	for rowIndex, row := range matrix {
		for colIndex, col := range row {
			result[rowIndex][colIndex] = depth * (factor*float64(col) - 0.5)
		}
	}
	return result
}
