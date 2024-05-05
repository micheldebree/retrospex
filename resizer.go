package main

import (
	"golang.org/x/image/draw"
	"image"
)

func Resize(img *image.Image, width, height int) image.Image {
	inRect := wholeImageRectangle(img)
	outRect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
	result := image.NewRGBA(outRect)
	draw.BiLinear.Scale(result, outRect, *img, inRect, draw.Over, nil)
	return result
}

// Rectangle covering the whole image
func wholeImageRectangle(img *image.Image) image.Rectangle {
	width := (*img).Bounds().Max.X - (*img).Bounds().Min.X
	height := (*img).Bounds().Max.Y - (*img).Bounds().Min.Y
	return image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
}
