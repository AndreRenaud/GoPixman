package pixman

import (
	"image"
	"image/color"
)

//var _ draw.Image = (*Image)(nil)

func ColorModelToPixman(c color.Model) PixmanFormatCode {
	switch c {
	case color.RGBAModel:
		return PIXMAN_a8r8g8b8
	case color.NRGBAModel:
		return PIXMAN_x8r8g8b8
	case color.NRGBAModel:
		return PIXMAN_x8b8g8r8
	default:
		return 0 // Unknown format
	}
}

func (i *Image) ColorModel() color.Model {
	format := ImageGetFormat(i)
	switch format {
	case PIXMAN_a8r8g8b8, PIXMAN_x8r8g8b8, PIXMAN_a8b8g8r8, PIXMAN_x8b8g8r8:
		return color.RGBAModel
	default:
		// TODO: Handle other formats
		return color.NRGBA64Model
	}
}

func (i *Image) Bounds() image.Rectangle {
	width := ImageGetWidth(i)
	height := ImageGetHeight(i)
	return image.Rect(0, 0, int(width), int(height))
}

//func (i *Image) At(x, y int) color.Color {
//data := ImageGetBits(i)
//stride := ImageGetStride(i)
//depth := ImageGetDepth(i)
