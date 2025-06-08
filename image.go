package pixman

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"unsafe"
)

var _ draw.Image = (*Image)(nil)

func (i *Image) ColorModel() color.Model {
	format := ImageGetFormat(i.pixman)
	switch format {
	case PIXMAN_a8r8g8b8, PIXMAN_x8r8g8b8, PIXMAN_a8b8g8r8, PIXMAN_x8b8g8r8:
		return color.RGBAModel
	case PIXMAN_b5g6r5, PIXMAN_r5g6b5:
		return color.RGBAModel // 16-bit formats, treated as RGB
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
	//case PIXMAN_b5g6r5:
	//col = color.RGBA{
	//R: (rawData[offset+1] & 0x1F) << 3,                            // Red bits
	//G: ((rawData[offset+1] >> 5) | (rawData[offset] & 0x03)) << 2, // Green bits
	//B: (rawData[offset] >> 2) & 0x1F << 3,                         // Blue bits
	//A: 255,                                                        // No alpha in this format
	//}
	case PIXMAN_r5g6b5:
		//log.Printf("RGB565 at %d,%d: %02x %02x", x, y, rawData[offset], rawData[offset+1])
		rgba := color.RGBA{
			R: (rawData[offset+1] & 0xf8),
			G: ((rawData[offset+1] & 0x7) << 5) | (rawData[offset]&0xe0)>>3,
			B: (rawData[offset] & 0x1f) << 3,
			A: 255, // No alpha in this format
		}
		if rgba.R&0x08 != 0 {
			rgba.R |= 0x07
		}
		if rgba.G&0x04 != 0 {
			rgba.G |= 0x03
		}
		if rgba.B&0x08 != 0 {
			rgba.B |= 0x07
		}
		col = rgba
		//log.Printf("RGB565 color at %d,%d: %v", x, y, col)
	default:
		col = color.Transparent // Unsupported format
	}
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

	colR, colG, colB, colA := c.RGBA()
	colR8 := uint8(colR >> 8)
	colG8 := uint8(colG >> 8)
	colB8 := uint8(colB >> 8)
	colA8 := uint8(colA >> 8)

	switch ImageGetFormat(i.pixman) {
	case PIXMAN_r8g8b8a8:
		rawData[offset] = colR8
		rawData[offset+1] = colG8
		rawData[offset+2] = colB8
		rawData[offset+3] = colA8
	case PIXMAN_r8g8b8x8:
		rawData[offset] = colR8
		rawData[offset+1] = colG8
		rawData[offset+2] = colB8
		rawData[offset+3] = 255
	//case PIXMAN_b5g6r5:
	//rawData[offset] = colB8&0xf8 | colG8>>5
	//rawData[offset+1] = colG8&0x1c<<3 | colR8>>3
	//case PIXMAN_r5g6b5:
	//rawData[offset] = colR8&0xf8 | (colG8>>5)&0x07
	//rawData[offset+1] = colG8&0x1c<<3 | colB8>>3
	default:
		log.Printf("Unsupported format for Set: %s", ImageGetFormat(i.pixman))
		// Unsupported format, do nothing
	}
}

// Composite performs a blit operation from the sub-image of `src` defined by `r`, placing the result at the point `sp` in this image.
func (i *Image) Composite(src *Image, r image.Rectangle, sp image.Point) {
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
