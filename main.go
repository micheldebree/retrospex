package main

// https://en.wikipedia.org/wiki/Locality-sensitive_hashing
// k means clustering

// https://stackoverflow.com/questions/41099138/k-means-versus-lsh-algorithm
// https://terenceshin.medium.com/top-five-clustering-algorithms-you-should-know-instead-of-k-means-clustering-b22f25e5bfb4

// https://stackoverflow.com/questions/37271413/heuristics-to-sort-array-of-2d-3d-points-according-their-mutual-distance/37308369#37308369
import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"golang.org/x/exp/maps"
)

var Version = "0.0"
var Arch = "dev"

const defaultOutput = "out.png"
const defaultMode = "koala"
const defaultPalette = "colodore"
const defaultDitherMatrix = "bayer4x4"
const defaultDitherDepth = 25

type Options struct {
	OutFile      string
	Mode         string
	Palette      string
	DitherMatrix string
	DitherDepth  int
}

func main() {

	var options Options

	flag.StringVar(&options.OutFile, "o", defaultOutput, "output filename")
	flag.StringVar(&options.Mode, "m", defaultMode, "graphics mode")
	flag.StringVar(&options.Palette, "p", defaultPalette, "palette")
	flag.StringVar(&options.DitherMatrix, "dm", defaultDitherMatrix, "dither matrix")
	flag.IntVar(&options.DitherDepth, "dd", defaultDitherDepth, "dither depth")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		help()
		return
	}

	spec, isPresent := C64Specs[options.Mode]
	if !isPresent {
		printError(fmt.Sprintf("Unknown mode: %s", options.Mode))
		return
	}

	palette, isPresent := C64Palettes[options.Palette]
	if !isPresent {
		printError(fmt.Sprintf("Unknown palette: %s", options.Palette))
		return
	}

	ditherMatrix, isPresent := DitherMatrices[options.DitherMatrix]
	if !isPresent {
		printError(fmt.Sprintf("Unknown dither matrix: %s", options.DitherMatrix))
		return
	}

	if options.DitherDepth < 0 || options.DitherDepth > 255 {
		printError(fmt.Sprintf("Unsupported dither depth: %d, must be 0-100", options.DitherDepth))
		return
	}

	infile := args[0]
	img, err := ReadImageFile(infile)
	if err != nil {
		printError(err.Error())
		return
	}
	img = CropFill(&img, spec)

	indexedImage := toIndexedImage(&img, spec, palette)
	OrderedDither(&indexedImage, ditherMatrix, options.DitherDepth)
	newImage := Quantize(indexedImage)

	result := newImage.Render()
	WriteImage(options.OutFile, result)
	fmt.Print(options.OutFile)
}

func printError(message string) {
	fmt.Print("\nERROR: ", message, "\n")
	help()
}

func help() {
	fmt.Printf("\nretrospex %s.%s by yth\n", Version, Arch)
	fmt.Printf("\nUsage: retrospex [options] input.png\n\n")
	fmt.Printf("Options:\n\n")
	fmt.Printf("\t-o\n\t\tOutput filename (default %s)\n", defaultOutput)
	fmt.Printf("\t-m\n\t\tGraphics mode. (default %s), One of %s\n", defaultMode, strings.Join(maps.Keys(C64Specs), ","))
	fmt.Printf("\t-p\n\t\tPalette (default %s). One of %s\n", defaultPalette, strings.Join(maps.Keys(C64Palettes), ","))
	fmt.Printf("\t-dm\n\t\tDither matrix (default %s). One of %s\n", defaultDitherMatrix, strings.Join(maps.Keys(DitherMatrices), ","))
	fmt.Printf("\t-dd\n\t\tDither depth (default %d). 0-255\n", defaultDitherDepth)
}
