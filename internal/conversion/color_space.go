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
	LINEAR_RGB ColorspaceName = "lineairRgb"
	XYZ        ColorspaceName = "xyz"
	XYY        ColorspaceName = "xyy"
	LAB        ColorspaceName = "lab"
)

func euclidianDistance(color1, color2 color.Color, colorspace ColorspaceName) float64 {

	ach1, ach2, ach3 := toColorSpace(toColorful(color1), colorspace)
	bch1, bch2, bch3 := toColorSpace(toColorful(color2), colorspace)

	diff1 := (bch1 - ach1)
	diff2 := (bch2 - ach2)
	diff3 := (bch3 - ach3)

	return math.Sqrt(diff1*diff1 + diff2*diff2 + diff3*diff3)
}

func toColorSpace(aColor colorful.Color, colorspace ColorspaceName) (float64, float64, float64) {

	switch colorspace {
	case RGB:
		return aColor.R, aColor.G, aColor.B
	case LINEAR_RGB:
		return aColor.LinearRgb()
	case XYZ:
		return aColor.Xyz()
	case XYY:
		return aColor.Xyy()
	case LAB:
		return aColor.Lab()
	default:
		panic(fmt.Sprintf("Colorspace %s is not supported", colorspace))
	}

}

func toColorful(aColor color.Color) colorful.Color {

	result, success := colorful.MakeColor(aColor)

	if !success {
		panic("Could not convert to colorspace")
	}

	return result

}
