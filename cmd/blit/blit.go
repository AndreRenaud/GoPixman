package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/AndreRenaud/go-pixman"
)

func loadFile(filename string) (image.Image, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	img, err := png.Decode(data)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	inputFile := flag.String("input", "", "Input image file")
	flag.Parse()

	img, err := loadFile(*inputFile)
	if err != nil {
		log.Fatalf("failed to load image %q: %v", *inputFile, err)
	}
	// Create a Pixman image from the RGBA image
	pixmanImage, err := pixman.ImageFromImage(img)
	if err != nil {
		log.Fatalf("failed to create Pixman image: %v", err)
	}
	// Get the format of the Pixman image
	format := pixman.ImageGetFormat(pixmanImage)
	log.Printf("Pixman image format: %s", format)
	pixSize := pixmanImage.Bounds()
	log.Printf("Pixman image size: %dx%d", pixSize.Dx(), pixSize.Dy())
}
