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
	"strconv"
	"strings"
	"time"

	"github.com/micheldebree/retrospex/internal/c64io"
	"github.com/micheldebree/retrospex/internal/conversion"
	"github.com/micheldebree/retrospex/internal/dithering"

	"github.com/micheldebree/retrospex/internal/imageio"
	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

var Version = "0.0"
var Arch = "dev"

type Options struct {
	OutFile          string
	Mode             string
	Palette          string
	DitherMatrix     string
	DitherDepth      int
	Format           FormatType // Changed to use enum type
	BitpatternColors string
}

type FormatType string

const (
	PNG FormatType = "png"
	BIN FormatType = "bin"
)

var defaultOptions = Options{
	OutFile:          "out.png",
	Mode:             "koala",
	Palette:          "colodore",
	DitherMatrix:     "bayer4x4",
	DitherDepth:      25,
	Format:           PNG, // Default format is png
	BitpatternColors: "",
}

func parseBitpatternColors(value string) map[int]int {

	if len(value) == 0 {
		return make(map[int]int, 0)
	}

	pairs := strings.Split(value, ",")
	result := make(map[int]int, len(pairs))

	for _, pair := range pairs {

		numbers := strings.Split(pair, ":")

		if len(numbers) != 2 {
			panic("Need two numbers to map bit pattern to color")
		}

		bitpattern, err := strconv.Atoi(numbers[0])
		if err != nil {
			panic("Illegal bitpattern")
		}
		color, err := strconv.Atoi(numbers[1])
		if err != nil {
			panic("Illegal color")
		}
		result[bitpattern] = color
	}
	return result
}

func main() {

	startTime := time.Now()

	fmt.Printf("\nretrospex %s.%s by yth\n", Version, Arch)

	var options Options

	flag.StringVar(&options.OutFile, "o", defaultOptions.OutFile, "Output filename")
	flag.StringVar(&options.Mode, "m", defaultOptions.Mode, "Graphics mode")
	flag.StringVar(&options.Palette, "p", defaultOptions.Palette, "Palette")
	flag.StringVar(&options.DitherMatrix, "dm", defaultOptions.DitherMatrix, "The name of a predefined ordered dithering matrix")
	flag.IntVar(&options.DitherDepth, "dd", defaultOptions.DitherDepth, "Dither depth (0-255). Depth of dithering.")
	flag.StringVar(&options.BitpatternColors, "bpc", defaultOptions.BitpatternColors, "Force bitpattern/color pairs. For example 0:0 to force background black.")

	var formatString string
	flag.StringVar(&formatString, "f", string(defaultOptions.Format), "Output format (png or binary)")
	flag.Parse()

	options.Format = FormatType(formatString)

	if options.Format != PNG && options.Format != BIN {
		printError(fmt.Sprintf("Unknown format: %s", options.Format))
		return
	}

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

	spec := indexedimage.MakeSpec(indexedimage.RetrospecName(options.Mode), &img)
	indexedImage := indexedimage.ToIndexedImage(&img, spec, palette)
	dithering.OrderedDither(&indexedImage, ditherMatrix, options.DitherDepth)

	forcedBitpatternMapping := parseBitpatternColors(options.BitpatternColors)

	newImage := conversion.Quantize(indexedImage, forcedBitpatternMapping)

	switch options.Format {
	case PNG:
		// TODO: add overwrite flag
		imageio.WriteImage(options.OutFile, newImage.Render())
	case BIN:
		c64io.SaveBinary(options.OutFile, &newImage, true)
	}

	fmt.Printf("%s --> %s (%s) in %s\n", infile, options.OutFile, options.Mode, time.Since(startTime))
}

func printError(message string) {
	fmt.Print("\nERROR: ", message, "\n")
	help()
}

func help() {
	fmt.Printf("\nUsage: retrospex [options] input.png\n\n")
	fmt.Printf("Options:\n\n")
	flag.PrintDefaults()
}
