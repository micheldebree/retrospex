package main

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// Quantize the center observation and use it as index for all pixels in the same cluster
func quantizeClusters(theClusters clusters.Clusters, palette Palette) []int {

	result := make([]int, len(theClusters))

	for i, c := range theClusters {
		aColor := toColorfulColor(c.Center)
		quantized, _ := QuantizeToIndex(aColor, palette)

		// convert Observations back to Pixel
		pixels := make([]Pixel, len(c.Observations))
		for ii, o := range c.Observations {
			pixels[ii] = o.(Pixel)
			pixels[ii].paletteIndex = quantized
			// we did not quantize all pixels in the cluster, so we don't know the quantization error
			pixels[ii].quantizationError = -1.0
		}

		fmt.Printf("Cluster: %d %+v (%+v) with %d pixels\n", i, c.Center, quantized, len(c.Observations))
		result[i] = quantized
	}
	return result
}

func (pixel Pixel) Coordinates() clusters.Coordinates {
	return clusters.Coordinates{pixel.color.R, pixel.color.G, pixel.color.B}
}

func (pixel Pixel) Distance(pos clusters.Coordinates) float64 {
	return distance(pixel.color, toColorfulColor(pos))
}

func toColorfulColor(pos clusters.Coordinates) colorful.Color {
	return colorful.Color{R: pos[0], G: pos[1], B: pos[2]}
}

// Cluster pixels
func cluster(pixels *[]Pixel, nrClusters int) clusters.Clusters {

	var d clusters.Observations

	for _, p := range *pixels {
		if !p.hasBitPattern() {
			d = append(d, p)
		}
	}

	km, err := kmeans.NewWithOptions(0.01, nil)
	abortOnError(err)

	result, err := km.Partition(d, nrClusters)
	abortOnError(err)

	return result
}

func reducePaletteKmeans(img IndexedImage, layer Layer) ReducedPalette {

	colorClusters := cluster(&img.pixels, len(layer.bitpatterns))
	quantizedMeans := quantizeClusters(colorClusters, img.palette)

	// TODO: make bitpattern map a type
	// TODO: does this make any difference?
	existingBitpatterns := make(map[int]int)
	for _, p := range img.pixels {
		if p.hasBitPattern() {
			existingBitpatterns[p.paletteIndex] = p.bitPattern
		}
	}

	newPalette := make(Palette)
	newBitpatterns := make(map[int]int)
	i := 0
	for _, quantizedMean := range quantizedMeans {
		newPalette[quantizedMean] = img.palette[quantizedMean]
		newBitpatterns[quantizedMean] = layer.bitpatterns[i]
		i++
	}
	// Add existing bitpatterns to the palette so they also get a chance
	// TODO: does this make any difference?
	for key := range existingBitpatterns {
		newPalette[key] = img.palette[key]
		newBitpatterns[key] = existingBitpatterns[key]
	}
	return ReducedPalette{newPalette, newBitpatterns}
}
