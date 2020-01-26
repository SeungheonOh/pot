package renderer

import (
	"fmt"
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

func ascii(imgOri cv.Mat) error {

	img := imgOri.Clone()
	cv.CvtColor(img, &img, 7)

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
