package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
				fmt.Printf("\033[48;2;%d;%d;%dm\033[30m▄\033[49m\033[39m",
					// Top Pixels
					imgPtr[i*img.Cols()*3+j+2],
					imgPtr[i*img.Cols()*3+j+1],
					imgPtr[i*img.Cols()*3+j])
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
	var ret image.Point
	if imgRatio > termRatio {
		ret = image.Point{X: int(float64(term.Y*2) / imgRatio), Y: term.Y * 2}
	} else {
		ret = image.Point{X: term.X, Y: int(float64(term.X) * imgRatio)}
	}
	return ret
}

func LoadFromURL(url string) (cv.Mat, error) {
	response, err := http.Get(url)
	if err != nil {
		return cv.NewMat(), errors.New("Failed load image")
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	img, err := cv.IMDecode(body, 1)

	if err != nil {
		return cv.NewMat(), errors.New("Failed decode image: ")
	}

	return img, nil
}

const (
	NAME        = "PixelOnTerminal"
	AUTHOR      = "SeungheonOh 2019|github.com/SeungheonOh/PixelOnTerminal"
	DESCRIPTION = "A cool guy way to view image/video/webcam in terminal enviroment"
	USAGE       = "pot <IMAGE/VIDEO/URL> <OPTIONS>"
)

var (
	fileName   string = ""
	fullScreen bool   = false
	dontExit   bool   = false
	useCam     bool   = false

	Arguments *flag.FlagSet
)

func init() {
	Arguments = flag.NewFlagSet("", flag.ExitOnError)

	Arguments.BoolVar(&fullScreen, "s", false, "Strach image to the size of terminal")
	Arguments.BoolVar(&fullScreen, "S", false, "Strach image to the size of terminal")

	Arguments.BoolVar(&dontExit, "e", false, "Repeat the input until keyboard interupt")
	Arguments.BoolVar(&dontExit, "E", false, "Repeat the input until keyboard interupt")

	Arguments.BoolVar(&useCam, "c", false, "Fetch webcam stream and print")
	Arguments.BoolVar(&useCam, "C", false, "Fetch webcam stream and print")
}

func main() {
	var fps int = 30

	var ArgIn []string

	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
		if len(arg) > 3 {
			fileName = arg // if argument is path
		} else {
			ArgIn = append(ArgIn, arg)
		}
	}

	Arguments.Parse(ArgIn)

	if fileName == "" && !useCam {
		// Help message
		fmt.Print(NAME, "\n")
		fmt.Print(DESCRIPTION, "\n\n")
		fmt.Print("USAGE:\n")
		fmt.Print("\t", USAGE, "\n\n")
		fmt.Print("OPTIONS:\n")
		Arguments.PrintDefaults()
		fmt.Print("\n", AUTHOR, "\n")
		return
	}

	img := cv.NewMat()
	defer img.Close()

	// Used for both image and video
	var capture *cv.VideoCapture
	if useCam {
		cam, err := cv.VideoCaptureDevice(0)
		fmt.Println("loading capture")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed open camera: %s\n", err.Error())
			return
		}
		capture = cam
	} else {
		video, err := cv.OpenVideoCapture(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed load image/video: %s\n", err.Error())
			return
		}
		capture = video
	}

	if !capture.IsOpened() && !useCam {
		img, err := LoadFromURL(fileName)
		if err != nil || img.Empty() {
			fmt.Fprintf(os.Stderr, "Failed load image: %s\n", err.Error())
			return
		}
	}

	fps = int(capture.Get(5)) // Get FPS value for video/GIF

	// Get initial terminal size
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))
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
		// Timer for precise frame calculation
		start := time.Now()

		ok := capture.Read(&img)
		if !ok {
			if dontExit {
				// Restart Video from beginning
				capture.Set(0, 0)
				continue
			}
			running = false
			return
		}

		var imgSize image.Point

		if fullScreen {
			imgSize = image.Point{X: width, Y: height * 2}
		} else {
			imgSize = CalculateSize(img, image.Point{X: width, Y: height})
		}

		cv.Resize(img, &img, imgSize, 0, 0, 1)

		err = PrintMat(img)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		end := time.Now()
		time.Sleep(time.Second/time.Duration(fps) - end.Sub(start))
	}
}
