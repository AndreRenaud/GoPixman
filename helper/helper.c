#include <pixman.h>
#include <stdio.h>

struct pixman_formats {
    int code;
    const char *name;
};
#define FORMAT(a) { a, #a }

#define FORMATS \
    FORMAT(PIXMAN_a8r8g8b8), \
    FORMAT(PIXMAN_x8r8g8b8), \
    FORMAT(PIXMAN_a8b8g8r8), \
    FORMAT(PIXMAN_x8b8g8r8), \
    FORMAT(PIXMAN_b8g8r8a8), \
    FORMAT(PIXMAN_b8g8r8x8), \
    FORMAT(PIXMAN_r5g6b5), \
    FORMAT(PIXMAN_b5g6r5), \
    FORMAT(PIXMAN_a1r5g5b5), \
    FORMAT(PIXMAN_x1r5g5b5), \
    FORMAT(PIXMAN_a1b5g5r5), \
    FORMAT(PIXMAN_x1b5g5r5), \
    FORMAT(PIXMAN_a4r4g4b4), \
    FORMAT(PIXMAN_x4r4g4b4), \
    FORMAT(PIXMAN_a4b4g4r4), \
    FORMAT(PIXMAN_x4b4g4r4)

int main(void)
{
    struct pixman_formats formats[] = {
        FORMATS
    };
    for (unsigned long i = 0; i < sizeof(formats) / sizeof(formats[0]); i++) {
        printf("%s PixmanFormatCode = 0x%08x\n", formats[i].name, formats[i].code);
    }

    return 0;
}
