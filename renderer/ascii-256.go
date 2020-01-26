package renderer

import (
	"errors"
	"fmt"
	"os"

	cv "gocv.io/x/gocv"
)

func init() {
	RendererMap["ascii-256"] = ascii256
}

func ascii256(img cv.Mat) error {
	imgPtr := img.DataPtrUint8()

	if img.Cols()*img.Rows()*3 != len(imgPtr) {
		return errors.New("Color RGB image is only supported")
	}

	fmt.Print("\033[0;0H")

	for i := 0; i < img.Rows(); i += 2 {
		for j := 0; j < img.Cols()*3; j += 3 {
			fmt.Fprintf(os.Stdout, "\033[38;2;%d;%d;%dm█\033[39m",
				imgPtr[i*img.Cols()*3+j+2],
				imgPtr[i*img.Cols()*3+j+1],
				imgPtr[i*img.Cols()*3+j])
		}
		fmt.Print("\n")
	}
	fmt.Print("\033[J")
	return nil
}
