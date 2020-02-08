package renderer

import (
	"errors"
	"image"

	cv "gocv.io/x/gocv"
)

var RendererMap = make(map[string]func(cv.Mat, image.Point) error)

type Renderer interface {
	Render(img cv.Mat) error
}

func Render(img cv.Mat, size image.Point, name string) error {
	renderer, exist := RendererMap[name]
	if !exist {
		return errors.New("Renderer Not Found")
	}
	return renderer(img, size)
}
