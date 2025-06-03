package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/AndreRenaud/GoPixman"
)

func loadFile(filename string) (image.Image, error) {
	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(data)
	if err != nil {
		_ = data.Close()
		return nil, err
	}
	if err := data.Close(); err != nil {
		return nil, err
	}
	return img, nil
}

func savePNG(img image.Image, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(outFile, img); err != nil {
		_ = outFile.Close()
		return err
	}
	if err := outFile.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	inputFile := flag.String("input", "", "Input image file")
	outputFile := flag.String("output", "output.png", "Output image file")
	rawFile := flag.String("raw", "output.raw", "Output raw image data file")
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
	pixSize := pixmanImage.Bounds()
	log.Printf("Pixman image size: %dx%d@%d", pixSize.Dx(), pixSize.Dy(), pixmanImage.Depth())

	//solid, err := pixman.ImageSolid(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	solid, err := pixman.ImageFromImage(img)
	//solid, err := pixman.ImageFromImage(image.NewUniform(color.RGBA{255, 0, 0, 255}))
	//draw.Draw(pixmanImage, image.Rect(0, 0, 20, 20), image.NewUniform(color.RGBA{255, 0, 0, 255}), image.Point{}, draw.Src)
	if err != nil {
		log.Fatalf("failed to create solid fill image: %s", err)
	}

	// Draw using Go's image/draw package
	draw.Draw(pixmanImage, image.Rect(0, 0, 20, 20), image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	// Fill a colour using pixman
	pixmanImage.Fill(image.Rect(10, 40, 5, 30), color.RGBA{128, 0, 128, 255})
	// Composite the images together using pixman
	pixmanImage.Composite(solid, image.Rect(10, 10, 300, 300), image.Pt(30, 30))

	if *outputFile != "" {
		if err := savePNG(pixmanImage, *outputFile); err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
		log.Printf("Saved %dx%d image to %q", pixSize.Dx(), pixSize.Dy(), *outputFile)
	}

	if *rawFile != "" {
		if err := pixmanImage.SaveRaw(*rawFile); err != nil {
			log.Fatalf("failed to save raw image: %v", err)
		}
		log.Printf("Saved raw image data to %q", *rawFile)
	}
}
