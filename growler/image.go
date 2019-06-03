package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
)

// ReadImage is used to read an image file with given name,
// it will return an error if found or will print the image pixels.
func ReadImage(filename string) ([][]Pixel, error) {
	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(filename)

	if err != nil {

		return nil, err
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		return nil, err
	}

	fmt.Printf("File: %s \nWidth:%d \nHeight: %d \n", filename, len(pixels[0]), len(pixels))

	return pixels, nil
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}
