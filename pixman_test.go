package pixman

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
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

func BenchmarkImageFill(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))
	col := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for i := 0; i < b.N; i++ {
		draw.Draw(img, img.Bounds(), &image.Uniform{C: col}, image.Point{}, draw.Src)
	}
}
func BenchmarkPixmanFill(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))
	pixmanImg, err := ImageFromImage(img)
	if err != nil {
		b.Fatalf("failed to create Pixman image: %v", err)
	}
	col := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for i := 0; i < b.N; i++ {
		pixmanImg.Fill(img.Bounds(), col)
	}
}

func BenchmarkImageBlit(b *testing.B) {
	img, err := loadFile("testdata/pg-coral.png")
	if err != nil {
		b.Fatalf("failed to load image: %v", err)
	}
	dest := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for i := 0; i < b.N; i++ {
		draw.Draw(dest, dest.Bounds(), img, image.Point{}, draw.Src)
	}
}

func BenchmarkPixmanBlit(b *testing.B) {
	img, err := loadFile("testdata/pg-coral.png")
	if err != nil {
		b.Fatalf("failed to load image: %v", err)
	}
	pixmanImg, err := ImageFromImage(img)
	if err != nil {
		b.Fatalf("failed to create Pixman image: %v", err)
	}
	dest := image.NewRGBA(image.Rect(0, 0, 320, 240))
	pixmanDest, err := ImageFromImage(dest)
	if err != nil {
		b.Fatalf("failed to create Pixman destination")
	}
	for i := 0; i < b.N; i++ {
		pixmanDest.Composite(pixmanImg, img.Bounds(), image.Point{X: 0, Y: 0})
	}
}
