package renderer

import (
	"fmt"
	"image"

	cv "gocv.io/x/gocv"
)

func init() {
	RendererMap["dummy"] = dummy
}

func dummy(img cv.Mat, size image.Point) error {
	fmt.Println("This is debug renderer for retriving image informations.")
	fmt.Println("Image Size:")
	fmt.Println("Cols: ", img.Cols(), " Rows: ", img.Rows())

	return nil
}
