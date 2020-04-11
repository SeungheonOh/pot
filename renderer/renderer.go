package renderer

import (
	"errors"
	"image"
  "fmt"
  "os"

	cv "gocv.io/x/gocv"
)

var RendererMap = make(map[string]func(cv.Mat, image.Point) (string, error))

func Render(img cv.Mat, size image.Point, name string) error {
	renderer, exist := RendererMap[name]
	if !exist {
		return errors.New("Renderer Not Found")
	}
  out, err := renderer(img, size)
	fmt.Fprintf(os.Stdout, "\033[0;0H")
  fmt.Fprintf(os.Stdout, out)
	// Clear from the end of the picture to the bottom of the tty
	// Also avoids leftover artifacts when image doesn't fill the tty
	fmt.Fprintf(os.Stdout, "\033[J")
  return err
}
