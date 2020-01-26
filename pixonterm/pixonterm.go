package pixonterm

import (
	//"errors"
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"

	"github.com/SeungheonOh/PixelOnTerminal/renderer"
	cv "gocv.io/x/gocv"
	"golang.org/x/crypto/ssh/terminal"
)

func SetTerm() {
	fmt.Print("\033c\033[?25l")
}

func RecoverTerm() {
	defer fmt.Print("\033[?25h\n")
}

func TermSize() image.Point {
	width, height, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return image.Point{width, height}
}

func PrintMat(img cv.Mat, rendererName string) error {
	return renderer.Render(img, rendererName)
}

func CalculateSize(img cv.Mat, term image.Point) image.Point {
	// This calculates the appropriate size for the frame, maintaining
	// the same aspect ratio, and filling the terminal window fully

	var termRatio = float64(term.Y*2) / float64(term.X)
	var imgRatio = float64(img.Rows()) / float64(img.Cols())
	var ret image.Point
	if imgRatio > termRatio {
		ret = image.Point{X: int(float64(term.Y*2) / imgRatio), Y: term.Y * 2}
	} else {
		ret = image.Point{X: term.X, Y: int(float64(term.X) * imgRatio)}
	}
	return ret
}

func EventHandler(running *bool, termsize *image.Point) {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)
	for {
		sig := <-signals
		switch sig {
		case os.Interrupt:
			*running = false
			RecoverTerm()
			close(signals)
			return
		case syscall.SIGWINCH:
			width, height, _ := terminal.GetSize(int(os.Stdout.Fd()))
			*termsize = image.Point{width, height}
		}
	}
}
