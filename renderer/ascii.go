package renderer

import (
	"fmt"
	"image"
	"os"

	cv "gocv.io/x/gocv"
)

/*
 Simple Brightness based Ascii renderer.
 No facny pixel masking or stuff
*/

var Bitmask = [...]rune{
	' ',
	' ',
	'/',
	'^',
	'*',
	'$',
	'@',
	'#',
}

func init() {
	RendererMap["ascii"] = ascii
}

func ascii(imgOri cv.Mat, size image.Point) error {
	cv.Resize(imgOri, &imgOri, size, 0, 0, 1)

	img := imgOri.Clone()
	cv.CvtColor(img, &img, 7)
	cv.GaussianBlur(img, &img, image.Point{3, 3}, 2, 2, 0)
	//cv.Canny(img, &img, 30, 50)
	//cv.Laplacian(img, &img, 10, 3, 1, 0, 0)
	//cv.Sobel(img, &img, 10, 1, 0, 3, 1, 0, 0)

	imgPtr := img.DataPtrUint8()

	fmt.Print("\033[0;0H")

	for i := 0; i < img.Rows(); i += 2 {
		for j := 0; j < img.Cols(); j++ {
			fmt.Fprintf(os.Stdout, string(Bitmask[(imgPtr[i*img.Cols()+j]>>5)]))
		}
		if i != img.Rows()-2 {
			fmt.Println("\033[K")
		}
	}

	fmt.Print("\033[J")
	return nil
}
