package pixman

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	pixmanLib uintptr

	ImageCreateBits      func(format PixmanFormatCode, width int32, height int32, bits *uint32, rowstride int32) *Image
	ImageCreateSolidFill func(color *PixmanColor) *Image
	ImageGetFormat       func(image *Image) PixmanFormatCode
	ImageGetWidth        func(image *Image) int32
	ImageGetHeight       func(image *Image) int32
	ImageGetStride       func(image *Image) int32
	ImageGetDepth        func(image *Image) int32
	ImageGetData         func(image *Image) *uint32
	ImageComposite       func(op int32, src *Image, mask *Image, dest *Image, srcX int32, srcY int32, maskX int32, maskY int32, destX int32, destY int32, width int32, height int32) uint32
)

type Image struct {
}

func init() {
	var err error
	pixmanLib, err = purego.Dlopen("/opt/homebrew/lib/libpixman-1.dylib", purego.RTLD_LAZY)
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
	purego.RegisterLibFunc(&ImageComposite, pixmanLib, "pixman_image_composite")
}

func ImageFromImage(img image.Image) (*Image, error) {
	bounds := img.Bounds()
	var format PixmanFormatCode
	var stride int32
	var bits *uint32
	switch t := img.(type) {
	case *image.RGBA:
	case *image.NRGBA:
		format = PIXMAN_a8r8g8b8
		stride = int32(t.Stride)
		bits = (*uint32)(unsafe.Pointer(&t.Pix[0]))
	default:
		return nil, fmt.Errorf("unsupported image format %T", img)
	}
	width := int32(bounds.Dx())
	height := int32(bounds.Dy())
	pixImg := ImageCreateBits(format, width, height, bits, stride)
	if pixImg == nil {
		return nil, fmt.Errorf("failed to create Pixman image")
	}
	return pixImg, nil
}
