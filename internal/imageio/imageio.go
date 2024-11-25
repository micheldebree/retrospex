package imageio

import (
	"image"
	"image/png"
	"os"
)

func abortOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// ReadImageFile read an image from a file
func ReadImageFile(filename string) (image.Image, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		abortOnError(f.Close())
	}(f)

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func WriteImage(filename string, image image.Image) {

	out, err := os.Create(filename)
	abortOnError(err)

	err = png.Encode(out, image)
	abortOnError(err)
}
