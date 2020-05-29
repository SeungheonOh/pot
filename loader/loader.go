package loader

import (
	"image"
)

var LoaderMap = make(map[string]func(...string) MediaLoader)

type MediaLoader interface {
	Load(filename string, size image.Point) ([]image.Image, error)
	ImageSize(filename string) (image.Point, error)
}
