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

func colorMatch(c1, c2 color.Color, delta uint32) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	rMatch := (r1>>8)+delta >= (r2>>8) && (r1>>8) <= (r2>>8)+delta
	gMatch := (g1>>8)+delta >= (g2>>8) && (g1>>8) <= (g2>>8)+delta
	bMatch := (b1>>8)+delta >= (b2>>8) && (b1>>8) <= (b2>>8)+delta
	aMatch := (a1>>8)+delta >= (a2>>8) && (a1>>8) <= (a2>>8)+delta
	return rMatch && gMatch && bMatch && aMatch
}

func compareSubImage(img1, img2 image.Image, bounds image.Rectangle, delta uint32) error {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			if !colorMatch(c1, c2, delta) {
				return fmt.Errorf("Pixel at (%d,%d) differs: img1=%#v, img2=%#v", x, y, c1, c2)
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

func savePng(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to create PNG file %s: %v", filename, err))
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(fmt.Sprintf("failed to encode PNG file %s: %v", filename, err))
	}
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

	if err := compareSubImage(pixmanImg, uniform, img.Bounds(), 0); err != nil {
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

	if err := compareSubImage(pixmanImg, img, img.Bounds(), 0); err != nil {
		t.Errorf("Image blit did not match expected image: %v", err)
	}
}

func TestRGB565(t *testing.T) {
	bases := []string{"testdata/red", "testdata/blue", "testdata/green", "testdata/pg-coral"}
	for _, base := range bases {
		t.Logf("Testing RGB565 for base: %s", base)
		img, err := loadFile(base + ".png")
		if err != nil {
			t.Fatalf("failed to load image: %v", err)
		}
		rawData, err := os.ReadFile(base + "-rgb565.raw")
		if err != nil {
			t.Fatalf("failed to read raw data file: %v", err)
		}
		bounds := img.Bounds()
		rawImage, err := ImageFromBits(PIXMAN_r5g6b5, bounds.Dx(), bounds.Dy(), rawData, bounds.Dx()*2)
		if err != nil {
			t.Fatalf("failed to create Pixman image from raw data: %v", err)
		}
		savePng(rawImage, base+"-rgb565.png")

		if err := compareSubImage(rawImage, img, rawImage.Bounds(), 0x7f); err != nil {
			t.Errorf("RGB565 image did not match expected image: %v", err)
		}
	}
}
