package renderer

import (
	"fmt"
	"image"
  "bytes"

	cv "gocv.io/x/gocv"
)

func init() {
	RendererMap["dummy"] = Dummy
}

func Dummy(img cv.Mat, size image.Point) (string, error) {
  var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "This is debug renderer for retriving image informations.")
	fmt.Fprintf(&buffer, "Image Size:")
  fmt.Fprintf(&buffer, "Cols: %d \nRows: %d", img.Cols(), img.Rows())

	return buffer.String(), nil
}
