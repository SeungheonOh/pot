package main

import (
	//"bufio"
	"fmt"
	cv "gocv.io/x/gocv"
	"image"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func PrintMat(img cv.Mat) {
	imgPtr := img.DataPtrUint8()

	for i := 0; i < img.Rows(); i += 2 {
		for j := 0; j < img.Cols()*3; j += 3 {
			if i+1 >= img.Rows() { // if height is not even, prevent to check pixel that isn't exist
				fmt.Printf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄\033[49m\033[39m",
					// Top Pixels
					imgPtr[i*img.Cols()*3+j+2],
					imgPtr[i*img.Cols()*3+j+1],
					imgPtr[i*img.Cols()*3+j],
					// Bottom Pixels
					0,
					0,
					0)
			} else {
				fmt.Printf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄\033[49m\033[39m",
					// Top Pixels
					imgPtr[i*img.Cols()*3+j+2],
					imgPtr[i*img.Cols()*3+j+1],
					imgPtr[i*img.Cols()*3+j],
					// Bottom Pixels
					imgPtr[(i+1)*img.Cols()*3+j+2],
					imgPtr[(i+1)*img.Cols()*3+j+1],
					imgPtr[(i+1)*img.Cols()*3+j])
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	webcam, _ := cv.VideoCaptureDevice(0)
	defer webcam.Close()

	img := cv.NewMat()
	defer img.Close()

	// For auto size window
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	termSize := strings.Fields(string(out))

	fmt.Print("\033[s\033c")
	for {
		fmt.Print("\033[u")
		webcam.Read(&img)

		Width, _ := strconv.Atoi(termSize[1])
		Height := Width * img.Rows() / img.Cols()
		cv.Resize(img, &img, image.Point{Width, Height}, 0, 0, 1)

		PrintMat(img)
		// Note :
		// if anyone know to to flush in golang, plase add flushing, that will make it flicker-free(I mean that moving cursor)
	}
}
