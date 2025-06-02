package pixman

type PixmanFormatCode uint32

// Pixman format codes (partial list, add more as needed)
// See https://gitlab.freedesktop.org/pixman/pixman/-/blob/9879f6cfc40b4ef3bdca4ee9aaedacff8fb87244/pixman/pixman.h#L1044
// Note: The lack of macros in Go means we have to manually define these
// See helper/helper.c to regenerate
const (
	PIXMAN_a8r8g8b8 PixmanFormatCode = 0x20028888
	PIXMAN_x8r8g8b8 PixmanFormatCode = 0x20020888
	PIXMAN_a8b8g8r8 PixmanFormatCode = 0x20038888
	PIXMAN_x8b8g8r8 PixmanFormatCode = 0x20030888
	PIXMAN_b8g8r8a8 PixmanFormatCode = 0x20088888
	PIXMAN_b8g8r8x8 PixmanFormatCode = 0x20080888
	PIXMAN_r5g6b5   PixmanFormatCode = 0x10020565
	PIXMAN_b5g6r5   PixmanFormatCode = 0x10030565
	PIXMAN_a1r5g5b5 PixmanFormatCode = 0x10021555
	PIXMAN_x1r5g5b5 PixmanFormatCode = 0x10020555
	PIXMAN_a1b5g5r5 PixmanFormatCode = 0x10031555
	PIXMAN_x1b5g5r5 PixmanFormatCode = 0x10030555
	PIXMAN_a4r4g4b4 PixmanFormatCode = 0x10024444
	PIXMAN_x4r4g4b4 PixmanFormatCode = 0x10020444
	PIXMAN_a4b4g4r4 PixmanFormatCode = 0x10034444
	PIXMAN_x4b4g4r4 PixmanFormatCode = 0x10030444
)

// PixmanColor mirrors the C struct pixman_color_t
// See: https://gitlab.freedesktop.org/pixman/pixman/-/blob/9879f6cfc40b4ef3bdca4ee9aaedacff8fb87244/pixman/pixman.h#L150
type PixmanColor struct {
	Red   uint16
	Green uint16
	Blue  uint16
	Alpha uint16
}
