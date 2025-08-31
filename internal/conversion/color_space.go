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

	switch colorspace {
	case RGB:
		return colorfulColor.R, colorfulColor.G, colorfulColor.B
	case LINEAR_RGB:
		return colorfulColor.LinearRgb()
	case XYZ:
		return colorfulColor.Xyz()
	case XYY:
		return colorfulColor.Xyy()
	case LAB:
		return colorfulColor.Lab()
	default:
		panic(fmt.Sprintf("Colorspace %s is not supported", colorspace))
	}

}
