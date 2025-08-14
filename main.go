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

	"github.com/micheldebree/retrospex/internal/c64io"
	"github.com/micheldebree/retrospex/internal/conversion"
	"github.com/micheldebree/retrospex/internal/dithering"

	"github.com/micheldebree/retrospex/internal/imageio"
	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"

	"golang.org/x/exp/maps"
)

var Version = "0.0"
var Arch = "dev"

type Options struct {
	OutFile      string
	Mode         string
	Palette      string
	DitherMatrix string
	DitherDepth  int
	Format       FormatType // Changed to use enum type
}

type FormatType string

const (
	PNG    FormatType = "png"
	BINARY FormatType = "binary"
)

var defaultOptions = Options{
	OutFile:      "out.png",
	Mode:         "koala",
	Palette:      "colodore",
	DitherMatrix: "bayer4x4",
	DitherDepth:  25,
	Format:       PNG, // Default format is png
}

func main() {

	var options Options

	flag.StringVar(&options.OutFile, "o", defaultOptions.OutFile, "output filename")
	flag.StringVar(&options.Mode, "m", defaultOptions.Mode, "graphics mode")
	flag.StringVar(&options.Palette, "p", defaultOptions.Palette, "palette")
	flag.StringVar(&options.DitherMatrix, "dm", defaultOptions.DitherMatrix, "dither matrix")
	flag.IntVar(&options.DitherDepth, "dd", defaultOptions.DitherDepth, "dither depth")

	var formatString string
	flag.StringVar(&formatString, "f", string(defaultOptions.Format), "output format (png or binary)")
	flag.Parse()

	options.Format = FormatType(formatString)

	args := flag.Args()

	if len(args) != 1 {
		help()
		return
	}

	palette, isPresent := pixels.C64Palettes[options.Palette]
	if !isPresent {
		printError(fmt.Sprintf("Unknown palette: %s", options.Palette))
		return
	}

	ditherMatrix, isPresent := dithering.DitherMatrices[options.DitherMatrix]
	if !isPresent {
		printError(fmt.Sprintf("Unknown dither matrix: %s", options.DitherMatrix))
		return
	}

	if options.DitherDepth < 0 || options.DitherDepth > 255 {
		printError(fmt.Sprintf("Unsupported dither depth: %d, must be 0-255", options.DitherDepth))
		return
	}

	infile := args[0]
	img, err := imageio.ReadImageFile(infile)
	if err != nil {
		printError(err.Error())
		return
	}

	spec := indexedimage.MakeSpec(options.Mode, &img)
	fmt.Printf("Mode: %s\n", options.Mode)
	indexedImage := indexedimage.ToIndexedImage(&img, spec, palette)
	dithering.OrderedDither(&indexedImage, ditherMatrix, options.DitherDepth)
	newImage := conversion.Quantize(indexedImage)

	switch options.Format {
	case PNG:
		// TODO: add overwrite flag
		imageio.WriteImage(options.OutFile, newImage.Render())
	case BINARY:
		c64io.SaveBinary(options.OutFile, &newImage, true)
	}
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
	fmt.Printf("\t-o\n\t\tOutput filename (default %s)\n", defaultOptions.OutFile)
	fmt.Printf("\t-m\n\t\tGraphics mode. (default %s), One of %s\n", defaultOptions.Mode, strings.Join(maps.Keys(indexedimage.RetrospecFactories), ","))
	fmt.Printf("\t-p\n\t\tPalette (default %s). One of %s\n", defaultOptions.Palette, strings.Join(maps.Keys(pixels.C64Palettes), ","))
	fmt.Printf("\t-dm\n\t\tDither matrix (default %s). One of %s\n", defaultOptions.DitherMatrix, strings.Join(maps.Keys(dithering.DitherMatrices), ","))
	fmt.Printf("\t-dd\n\t\tDither depth (default %d). 0-255\n", defaultOptions.DitherDepth)
	fmt.Printf("\t-f\n\t\tOutput format (default %s). One of png, binary\n", defaultOptions.Format) // Added Format option to help
}
