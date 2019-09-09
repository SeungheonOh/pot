package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"

	cv "gocv.io/x/gocv"
	"golang.org/x/crypto/ssh/terminal"
)

var running = true
var width, height int

func PrintMat(img cv.Mat) error {
	imgPtr := img.DataPtrUint8()
	if img.Cols()*img.Rows()*3 != len(imgPtr) {
		return errors.New("Only supports Color RGB image for now")
	}

	// Move cursor to 0, 0
	fmt.Print("\033[0;0H")

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
		// Don't draw a newline at the bottom of the terminal
		// Also clear all characters to the edge of the terminal
		// Prevents artifacts when the image.width < tty.width
		if i != img.Rows()-2 {
			fmt.Println("\033[K")
		}
	}
	// Clear from the end of the picture to the bottom of the tty
	// Also avoids leftover artifacts when image doesn't fill the tty
	fmt.Print("\033[J")
	return nil
}

func CalculateSize(img cv.Mat, term image.Point) image.Point {
	// This calculates the appropriate size for the frame, maintaining
	// the same aspect ratio, and filling the terminal window fully

	var termRatio = float64(term.Y*2) / float64(term.X)
	var imgRatio = float64(img.Rows()) / float64(img.Cols())
	if imgRatio > termRatio {
		return image.Point{X: int(float64(term.Y*2) / imgRatio), Y: term.Y * 2}
	} else {
		return image.Point{X: term.X, Y: int(float64(term.X) * imgRatio)}
	}
}

func main() {
	webcam, err := cv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open camera: %s\n", err.Error())
		return
	}
	defer webcam.Close()

	img := cv.NewMat()
	defer img.Close()

	// Get initial terminal size
	width, height, err = terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get terminal size: %s\n", err.Error())
		return
	}

	// Handle SIGINT and stop the loop cleanly
	// Handle SIGWINCH to get new terminal size
	go func() {
		var signals = make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)
		for {
			sig := <-signals
			switch sig {
			case os.Interrupt:
				running = false
				close(signals)
				return
			case syscall.SIGWINCH:
				width, height, _ = terminal.GetSize(int(os.Stdout.Fd()))
			}
		}
	}()

	// Be sure to enable cursor when we exit
	defer fmt.Print("\033[?25h")

	// Reset/clear terminal and hide cursor
	fmt.Print("\033c\033[?25l")

	for running {
		ok := webcam.Read(&img)
		if !ok {
			fmt.Fprintln(os.Stderr, "Failed to read from camera")
			return
		}

		cv.Resize(img, &img, CalculateSize(img, image.Point{X: width, Y: height}), 0, 0, 1)

		err = PrintMat(img)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}
}
