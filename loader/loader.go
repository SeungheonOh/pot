package loader

import (
	"image"
)

type MediaLoader interface {
	Load(filename string, size image.Point) ([]image.Image, error)
}
