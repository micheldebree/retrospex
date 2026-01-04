package conversion

import (
	"fmt"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type ColorspaceName string

const (
	RGB        ColorspaceName = "rgb"
	LINEAR_RGB ColorspaceName = "linearRgb"
	XYZ        ColorspaceName = "xyz"
	XYY        ColorspaceName = "xyy"
	LAB        ColorspaceName = "lab"
)

var ColorSpaceConverters = map[ColorspaceName]func(colorful.Color) (float64, float64, float64){
	RGB:        func(color colorful.Color) (float64, float64, float64) { return color.R, color.G, color.B },
	LINEAR_RGB: func(color colorful.Color) (float64, float64, float64) { return color.LinearRgb() },
	XYZ:        func(color colorful.Color) (float64, float64, float64) { return color.Xyz() },
	XYY:        func(color colorful.Color) (float64, float64, float64) { return color.Xyy() },
	LAB:        func(color colorful.Color) (float64, float64, float64) { return color.Lab() },
}

func euclidianDistance(color1, color2 color.Color, colorspace ColorspaceName) float64 {

	ach1, ach2, ach3 := toColorSpace(color1, colorspace)
	bch1, bch2, bch3 := toColorSpace(color2, colorspace)

	diff1 := (bch1 - ach1)
	diff2 := (bch2 - ach2)
	diff3 := (bch3 - ach3)

	return math.Sqrt(diff1*diff1 + diff2*diff2 + diff3*diff3)
}

func toColorSpace(aColor color.Color, colorspace ColorspaceName) (float64, float64, float64) {

	colorfulColor, success := colorful.MakeColor(aColor)

	if !success {
		panic("Could not convert to colorspace")
	}

	converter, isPresent := ColorSpaceConverters[colorspace]

	if !isPresent {
		panic(fmt.Sprintf("Colorspace %s is not supported", colorspace))
	}

	return converter(colorfulColor)

}
