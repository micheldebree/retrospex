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
	"github.com/micheldebree/retrospex/internal/prepost"
)

var Version = "0.0"
var Arch = "dev"

type Options struct {
	OutFile          string
	Mode             string
	Palette          string
	ColorSpace       string
	DitherMatrix     string
	DitherDepth      int
	Format           FormatType // Changed to use enum type
	BitpatternColors string
	AllowOverwrite   bool
	Resize           bool
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
	ColorSpace:       "linearRgb",
	DitherMatrix:     "bayer4x4",
	DitherDepth:      25,
	Format:           PNG,
	BitpatternColors: "",
	AllowOverwrite:   false,
	Resize:           false,
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
	flag.StringVar(&options.Mode, "m", defaultOptions.Mode, fmt.Sprintf("Graphics mode.\nOne of %s", KeysToString(indexedimage.RetrospecFactories)))
	flag.StringVar(&options.Palette, "p", defaultOptions.Palette, fmt.Sprintf("Palette.\nOne of %s", KeysToString(pixels.C64Palettes)))
	flag.StringVar(&options.ColorSpace, "cs", defaultOptions.ColorSpace, fmt.Sprintf("Colorspace conversion.\nOne of %s", KeysToString(conversion.ColorSpaceConverters)))
	flag.StringVar(&options.DitherMatrix, "dm", defaultOptions.DitherMatrix, fmt.Sprintf("A preset ordered dithering matrix.\nOne of %s", KeysToString(dithering.DitherMatrices)))
	flag.IntVar(&options.DitherDepth, "dd", defaultOptions.DitherDepth, "Dither depth (0-255). Depth of dithering.")
	flag.StringVar(&options.BitpatternColors, "bpc", defaultOptions.BitpatternColors, `Force bitpattern/color pairs. 
For example 0:0 to force background black.
For example 0:0,1:1,2:15,3:13 to force colors for all 4 bitpatterns`)
	flag.BoolVar(&options.Resize, "resize", defaultOptions.Resize, "Resize input image to fill target dimensions")
	flag.BoolVar(&options.AllowOverwrite, "overwrite", defaultOptions.AllowOverwrite, "Allow overwriting output file")

	var formatString string
	flag.StringVar(&formatString, "f", string(defaultOptions.Format), `Output format.
png for an image to convert for instance with png2prg
bin for binary data to use in Commodore64 development (see documentation for structure)`)
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

	// Get bitpattern size from hardcoded map to avoid chicken-egg problem
	bitPatternSize := indexedimage.GetBitPatternSize(indexedimage.RetrospecName(options.Mode))
	img = prepost.ResizeInputByPatternSize(img, bitPatternSize, options.Resize)
	spec := indexedimage.MakeSpec(indexedimage.RetrospecName(options.Mode), &img)
	indexedImage := indexedimage.ToIndexedImage(&img, spec, palette)
	dithering.OrderedDither(&indexedImage, ditherMatrix, options.DitherDepth)

	forcedBitpatternMapping := parseBitpatternColors(options.BitpatternColors)

	newImage := conversion.Quantize(indexedImage, forcedBitpatternMapping, conversion.ColorspaceName(options.ColorSpace))

	switch options.Format {
	case PNG:
		// TODO: add overwrite flag
		imageio.WriteImage(options.OutFile, newImage.Render(), options.AllowOverwrite)
	case BIN:
		c64io.SaveBinary(options.OutFile, &newImage, options.AllowOverwrite)
	}

	spec.Print()

	fmt.Printf("%s --> %s (%s) in %s\n", infile, options.OutFile, options.Mode, time.Since(startTime))
}

func KeysToString[K comparable, V any](m map[K]V) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, fmt.Sprintf("%v", k))
	}
	return strings.Join(keys, ", ")
}

func printError(message string) {
	fmt.Printf("\nERROR: %s\n", message)
	help()
}

func help() {
	fmt.Printf("\nUsage: retrospex [options] input.png\n\n")
	fmt.Printf("Options:\n\n")
	flag.PrintDefaults()
}
