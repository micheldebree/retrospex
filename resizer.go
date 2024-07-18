package main

import (
	"image"

	"golang.org/x/image/draw"
)

func CropFill(img *image.Image, spec Retrospec) image.Image {

	imgW := (*img).Bounds().Max.X - (*img).Bounds().Min.X
	imgH := (*img).Bounds().Max.Y - (*img).Bounds().Min.Y

	ratio := float64(imgW) / float64(imgH)
	newRatio := float64(spec.displayWidth()) / float64(spec.displayHeight())

	var newW, newH int
	if ratio > newRatio {
		newW = int(float64(imgH) * newRatio)
		newH = imgH
	} else {
		newW = imgW
		newH = int(float64(imgW) / newRatio)
	}

	cropLeft := (imgW - newW) / 2
	cropTop := (imgH - newH) / 2

	topLeft := image.Point{X: cropLeft, Y: cropTop}
	bottomRight := image.Point{X: cropLeft + newW, Y: cropTop + newH}

	srcRect := image.Rectangle{Min: topLeft, Max: bottomRight}
	destRect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: spec.width, Y: spec.height}}

	result := image.NewRGBA(destRect)

	draw.BiLinear.Scale(result, destRect, *img, srcRect, draw.Over, nil)
	return result
}
