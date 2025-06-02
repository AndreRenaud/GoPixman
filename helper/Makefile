CFLAGS=-Wall -Wextra -pedantic -std=c99 $(shell pkg-config --cflags pixman-1)
LFLAGS=$(shell pkg-config --libs pixman-1)
helper: helper.c
	gcc -o helper helper.c $(CFLAGS) $(LFLAGS)

clean:
	rm -f helper
