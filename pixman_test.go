package pixman

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func compareSubImage(img1, img2 image.Image, bounds image.Rectangle) error {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return fmt.Errorf("Pixel at (%d,%d) differs: img1=%v, img2=%v", x, y, c1, c2)
			}
		}
	}
	return nil
}

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
		return nil, fmt.Errorf("failed to close file %s: %v", filename, err)
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

func TestImageFill(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))
	col := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	pixmanImg, err := ImageFromImage(img)
	if err != nil {
		t.Fatalf("failed to create Pixman image: %v", err)
	}

	pixmanImg.Fill(img.Bounds(), col)
	uniform := &image.Uniform{C: col}

	if err := compareSubImage(pixmanImg, uniform, img.Bounds()); err != nil {
		t.Errorf("Image fill did not match expected color: %v", err)
	}
}

func TestImageBlit(t *testing.T) {
	img, err := loadFile("testdata/pg-coral.png")
	if err != nil {
		t.Fatalf("failed to load image: %v", err)
	}
	srcImage, err := ImageFromImage(img)
	if err != nil {
		t.Fatalf("failed to create Pixman image: %v", err)
	}

	dest := image.NewRGBA(img.Bounds())
	pixmanImg, err := ImageFromImage(dest)
	if err != nil {
		t.Fatalf("failed to create Pixman image: %v", err)
	}
	pixmanImg.Composite(srcImage, img.Bounds(), image.Point{X: 0, Y: 0})

	if err := compareSubImage(pixmanImg, img, img.Bounds()); err != nil {
		t.Errorf("Image blit did not match expected image: %v", err)
	}
}
