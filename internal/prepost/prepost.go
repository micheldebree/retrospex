package prepost

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/micheldebree/retrospex/internal/pixels"
)

// ResizeInputByPatternSize resizes input image using just the bitpattern size, avoiding chicken-egg problem
func ResizeInputByPatternSize(img image.Image, bitPatternSize int, shouldResize bool) image.Image {
	// For 2-bit modes (targetWidth == 160) we need to squash the picture horizontally
	// because each pixel is double width. Use a direct resize to preserve the aspect ratio
	// vertically while halving the horizontal dimension.

	originalWidth, originalHeight := pixels.GetDimensions(&img)
	targetWidth, targetHeight := 320/bitPatternSize, 200

	// adjust horizontal scale when double pixels
	if bitPatternSize > 1 {
		scaledWidth := originalWidth / bitPatternSize
		img = imaging.Resize(img, scaledWidth, originalHeight, imaging.Lanczos)
	}

	if shouldResize {
		return resizeImage(img, targetWidth, targetHeight)
	}
	return img
}

func resizeImage(img image.Image, targetWidth, targetHeight int) image.Image {
	// Scale to fill target dimensions, maintaining aspect ratio and crop excess
	return imaging.Fill(img, targetWidth, targetHeight, imaging.Center, imaging.NearestNeighbor)
}
