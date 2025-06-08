# GoPixman

[![Go Ci](https://github.com/AndreRenaud/GoPixman/actions/workflows/go.yml/badge.svg)](https://github.com/AndreRenaud/GoPixman/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AndreRenaud/GoPixman)](https://goreportcard.com/report/github.com/AndreRenaud/GoPixman)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/AndreRenaud/GoPixman)](https://pkg.go.dev/github.com/AndreRenaud/GoPixman)


GoPixman is a purego wrapper for Pixman to provide highly optimised software blitting functionality, especially on embedded platforms.

Pixman is a library that provides low-level pixel manipulation features such as image compositing and trapezoid rasterization.

Canonical Pixman source: https://gitlab.freedesktop.org/pixman/pixman

## Helpers
Use FFMPEG to create raw images for format testing, ie:
```
ffmpeg -i testdata/pg-coral.png -t 5 -r 1 -pix_fmt rgb565 frame-%d.raw
```