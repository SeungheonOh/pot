package main

import (
	"errors"
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
				fmt.Printf("\033[48;2;%d;%d;%dm\033[30m▄\033[49m\033[48m",
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
	if imgRatio > termRatio {
		return image.Point{X: int(float64(term.Y*2) / imgRatio), Y: term.Y * 2}
	} else {
		return image.Point{X: term.X, Y: int(float64(term.X) * imgRatio)}
	}
}

func LoadFromUrl(url string) (cv.Mat, error) {
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
	OPTIONS = [...]string{
		"-F",
		"-E",
		"-C",
	}
	OPTION_DESCRIPTION = [...]string{
		"Full screen mode, discard screen ratio",
		"Repeat Video, Gif, Image",
		"Load webcam, image/Video/URL is not required for this option",
	}
)

func main() {
	img := cv.NewMat()
	defer img.Close()

	var fps int = 30

	// Options
	var fileName string = ""
	var fullScreen bool = false
	var dontExit bool = false
	var useCam bool = false

	_ = useCam
	_ = fileName
	_ = dontExit
	_ = fullScreen

	for i := 0; i < len(os.Args); i++ {
		if os.Args[i][0] == '-' {
			fmt.Println(os.Args[i])
			// If argument is option
			switch os.Args[i] {
			case "-F", "-f":
				fullScreen = true
				break
			case "-E", "-e":
				dontExit = true
				break
			case "-C", "-c":
				useCam = true
				break
			default:
				break
			}
		} else if os.Args[i][len(os.Args[i])-4:][0] == '.' || os.Args[i][len(os.Args[i])-5:][0] == '.' {
			fileName = os.Args[i]
		}
	}

	if fileName == "" && !useCam {
		fmt.Print(NAME, "\n")
		fmt.Print(DESCRIPTION, "\n\n")
		fmt.Print("USAGE:\n")
		fmt.Print("\t", USAGE, "\n\n")
		fmt.Print("OPTIONS:\n")
		for i, option := range OPTIONS {
			fmt.Print("\t", option, "   ", OPTION_DESCRIPTION[i], "\n")
		}
		fmt.Print("\n", AUTHOR, "\n")
		return
	}

	var capture *cv.VideoCapture // Used for both image and video
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
		img, err := LoadFromUrl(fileName)
		if err != nil || img.Empty() {
			fmt.Fprintf(os.Stderr, "Failed load image: %s\n", err.Error())
			return
		}
	}

	fps = int(capture.Get(5))

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
		start := time.Now()
		ok := capture.Read(&img)
		if !ok {
			if dontExit {
				capture.Set(0, 0) // Restart Video from beginning
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
		//delayFor := int(int64() )
		time.Sleep(time.Second/time.Duration(fps) - end.Sub(start))
		_ = start
		_ = end
	}
}
