package pixman

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"unsafe"
)

var _ draw.Image = (*Image)(nil)

func (i *Image) ColorModel() color.Model {
	format := ImageGetFormat(i.pixman)
	switch format {
	case PIXMAN_a8r8g8b8, PIXMAN_x8r8g8b8, PIXMAN_a8b8g8r8, PIXMAN_x8b8g8r8:
		return color.RGBAModel
	default:
		// TODO: Handle other formats
		return color.NRGBA64Model
	}
}

func (i *Image) Bounds() image.Rectangle {
	width := ImageGetWidth(i.pixman)
	height := ImageGetHeight(i.pixman)
	return image.Rect(0, 0, int(width), int(height))
}

func (i *Image) Depth() int {
	depth := ImageGetDepth(i.pixman)
	if depth <= 0 {
		return 0
	}
	return int(depth)
}

func (i *Image) getRawData() []byte {
	return i.rawData
	/*
		dataPtr := ImageGetData(i.pixman)
		if dataPtr == nil {
			return nil
		}
		// Convert the pointer to a byte slice
		size := int(ImageGetStride(i.pixman) * ImageGetHeight(i.pixman))
		return unsafe.Slice((*byte)(unsafe.Pointer(dataPtr)), size)
	*/
}

func (i *Image) At(x, y int) color.Color {
	stride := ImageGetStride(i.pixman)
	depth := ImageGetDepth(i.pixman)
	if x < 0 || y < 0 || x >= int(ImageGetWidth(i.pixman)) || y >= int(ImageGetHeight(i.pixman)) {
		return color.Transparent
	}
	if stride <= 0 || depth <= 0 {
		return color.Transparent
	}
	offset := y*int(stride) + x*int(depth)/8
	rawData := i.getRawData()

	var col color.Color

	switch ImageGetFormat(i.pixman) {
	case PIXMAN_r8g8b8a8:
		col = color.RGBA{
			R: rawData[offset],
			G: rawData[offset+1],
			B: rawData[offset+2],
			A: rawData[offset+3],
		}
	case PIXMAN_r8g8b8x8:
		col = color.NRGBA{
			R: rawData[offset],
			G: rawData[offset+1],
			B: rawData[offset+2],
			A: 255,
		}
	default:
		col = color.Transparent // Unsupported format
	}
	//if (y%20) == 0 && (x%20) == 0 {
	//log.Printf("At(%d, %d) offset: %d/%d, stride: %d, depth: %d = %#v", x, y, offset, len(i.rawData), stride, depth, col)
	//}
	return col
}

func (i *Image) Set(x, y int, c color.Color) {
	stride := ImageGetStride(i.pixman)
	depth := ImageGetDepth(i.pixman)
	if x < 0 || y < 0 || x >= int(ImageGetWidth(i.pixman)) || y >= int(ImageGetHeight(i.pixman)) {
		return
	}
	if stride <= 0 || depth <= 0 {
		return
	}
	offset := y*int(stride) + x*int(depth)/8
	rawData := i.getRawData()

	switch ImageGetFormat(i.pixman) {
	case PIXMAN_r8g8b8a8:
		if rgba, ok := c.(color.RGBA); ok {
			rawData[offset] = rgba.R
			rawData[offset+1] = rgba.G
			rawData[offset+2] = rgba.B
			rawData[offset+3] = rgba.A
		}
	case PIXMAN_r8g8b8x8:
		if rgba, ok := c.(color.NRGBA); ok {
			rawData[offset] = rgba.R
			rawData[offset+1] = rgba.G
			rawData[offset+2] = rgba.B
			rawData[offset+3] = 255 // No alpha in this format
		}
	default:
		// Unsupported format, do nothing
	}
}

// Composite performs a blit operation from the sub-image of `src` defined by `r`, placing the result at the point `sp` in this image.
func (i *Image) Composite(src *Image, r image.Rectangle, sp image.Point) {
	//log.Printf("Blitting %d,%d-%d,%d to %d,%d-%d,%d", r.Min.X, r.Min.Y, r.Min.X+r.Dx(), r.Min.Y+r.Dy(), sp.X, sp.Y, sp.X+r.Dx(), sp.Y+r.Dy())
	//log.Printf("Src: %p %v, Dest: %p %v", src.pixman, src.Bounds(), i.pixman, i.Bounds())

	ImageComposite32(PIXMAN_OP_OVER, src.pixman, nil, i.pixman,
		int32(r.Min.X), int32(r.Min.Y), // src_x, src_y (source rectangle)
		0, 0, // mask_x, mask_y (no mask)
		int32(sp.X), int32(sp.Y), // dest_x, dest_y (destination point)
		int32(r.Dx()), int32(r.Dy())) // width, height (rectangle size)

}

func (i *Image) SaveRaw(filename string) error {
	if err := os.WriteFile(filename, i.getRawData(), 0644); err != nil {
		return err
	}
	return nil
}

func (i *Image) Fill(rect image.Rectangle, col color.Color) {
	rawData := i.getRawData()
	r, g, b, a := col.RGBA()
	col32 := uint32((a>>8)&0xff)<<24 | // Alpha
		uint32((b>>8)&0xff)<<16 | // Blue
		uint32((g>>8)&0xff)<<8 | // Green
		uint32((r>>8)&0xff) // Red
	stride := int(ImageGetStride(i.pixman) / 4) // Rowstride in 32-bit units
	Fill((*uint32)(unsafe.Pointer(&rawData[0])), stride, i.Depth(), rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(), col32)
}
