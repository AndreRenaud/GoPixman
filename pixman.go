package pixman

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	pixmanLib uintptr

	// These must match the C function signatures
	ImageCreateBits      func(format PixmanFormatCode, width int, height int, bits *uint32, rowstride int) *PixmanImage
	ImageCreateSolidFill func(color *PixmanColor) *PixmanImage
	ImageGetFormat       func(image *PixmanImage) PixmanFormatCode
	ImageGetWidth        func(image *PixmanImage) int32
	ImageGetHeight       func(image *PixmanImage) int32
	ImageGetStride       func(image *PixmanImage) int32
	ImageGetDepth        func(image *PixmanImage) int32
	ImageGetData         func(image *PixmanImage) *uint32
	ImageComposite32     func(op PixmanOperation, src *PixmanImage, mask *PixmanImage, dest *PixmanImage, src_x, src_y, mask_x, mask_y, dest_x, dest_y int32, width, height int32)
	Fill                 func(bits *uint32, stride int, bpp int, x int, y int, width int, height int, xor uint32) int
	ImageUnref           func(image *PixmanImage) int
)

type Image struct {
	rawData []byte
	pixman  *PixmanImage
}

type PixmanImage struct{}

func findPixmanLibrary() string {
	var dirs []string
	libraryName := "libpixman-1.so.0"
	switch runtime.GOOS {
	case "darwin":
		dirs = append(dirs, "/opt/homebrew/lib/", "/usr/local/lib", "/usr/lib")
		libraryName = "libpixman-1.dylib"
	case "linux":
		dirs = append(dirs, "/usr/lib", "/usr/lib/x86_64-linux-gnu", "/usr/local/lib") // TODO: Parse /etc/ld.so.conf
	case "windows":
		dirs = append(dirs, "C:\\Windows\\System32", "C:\\MinGW64\\bin")
		libraryName = "libpixman-1-0.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}

	for _, dir := range dirs {
		filename := dir + "/" + libraryName
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}
	panic(fmt.Sprintf("%s not found in %v", libraryName, dirs))
}

func init() {
	var err error
	pixmanLib, err = purego.Dlopen(findPixmanLibrary(), purego.RTLD_LAZY)
	if err != nil {
		panic("failed to load libpixman-1: " + err.Error())
	}

	purego.RegisterLibFunc(&ImageCreateBits, pixmanLib, "pixman_image_create_bits")
	purego.RegisterLibFunc(&ImageCreateSolidFill, pixmanLib, "pixman_image_create_solid_fill")
	purego.RegisterLibFunc(&ImageGetFormat, pixmanLib, "pixman_image_get_format")
	purego.RegisterLibFunc(&ImageGetWidth, pixmanLib, "pixman_image_get_width")
	purego.RegisterLibFunc(&ImageGetHeight, pixmanLib, "pixman_image_get_height")
	purego.RegisterLibFunc(&ImageGetStride, pixmanLib, "pixman_image_get_stride")
	purego.RegisterLibFunc(&ImageGetDepth, pixmanLib, "pixman_image_get_depth")
	purego.RegisterLibFunc(&ImageGetData, pixmanLib, "pixman_image_get_data")
	purego.RegisterLibFunc(&ImageComposite32, pixmanLib, "pixman_image_composite32")
	purego.RegisterLibFunc(&ImageUnref, pixmanLib, "pixman_image_unref")
	purego.RegisterLibFunc(&Fill, pixmanLib, "pixman_fill")
}

func ImageFromImage(img image.Image) (*Image, error) {
	// We don't do subimages yet
	if img.Bounds().Min.X != 0 || img.Bounds().Min.Y != 0 {
		return nil, fmt.Errorf("image bounds must start at (0,0), got %v", img.Bounds())
	}
	bounds := img.Bounds()
	var format PixmanFormatCode
	var stride int
	var bits *uint32
	switch t := img.(type) {
	case *image.RGBA:
		format = PIXMAN_r8g8b8a8
		stride = t.Stride
		bits = (*uint32)(unsafe.Pointer(&t.Pix[0]))
	case *image.NRGBA:
		format = PIXMAN_r8g8b8a8
		stride = t.Stride
		bits = (*uint32)(unsafe.Pointer(&t.Pix[0]))
	default:
		return nil, fmt.Errorf("unsupported image format %T", img)
	}
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: width=%d, height=%d", width, height)
	}
	bitSlice := unsafe.Slice((*uint8)(unsafe.Pointer(bits)), height*stride)
	retval := &Image{
		rawData: bitSlice,
	}
	retval.pixman = ImageCreateBits(format, width, height, bits, stride)

	if retval.pixman == nil {
		return nil, fmt.Errorf("failed to create Pixman image")
	}
	runtime.AddCleanup(retval, func(raw *PixmanImage) {
		ImageUnref(raw)
	}, retval.pixman)
	return retval, nil
}

func ImageSolid(col color.Color) (*Image, error) {
	r, g, b, a := col.RGBA()
	pixCol := &PixmanColor{
		Red:   uint16(r),
		Green: uint16(g),
		Blue:  uint16(b),
		Alpha: uint16(a),
	}
	retval := &Image{}
	retval.pixman = ImageCreateSolidFill(pixCol)
	if retval.pixman == nil {
		return nil, fmt.Errorf("failed to create Pixman solid fill image")
	}
	runtime.AddCleanup(retval, func(raw *PixmanImage) {
		ImageUnref(raw)
	}, retval.pixman)
	return retval, nil
}
