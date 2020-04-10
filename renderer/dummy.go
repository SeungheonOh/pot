package renderer

import (
	"fmt"
	"image"
  "bytes"

	cv "gocv.io/x/gocv"
)

func init() {
	RendererMap["dummy"] = dummy
}

func dummy(img cv.Mat, size image.Point) (string, error) {
  var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "This is debug renderer for retriving image informations.")
	fmt.Fprintf(&buffer, "Image Size:")
	fmt.Fprintf(&buffer, "Cols: ", img.Cols(), " Rows: ", img.Rows())

	return buffer.String(), nil
}
