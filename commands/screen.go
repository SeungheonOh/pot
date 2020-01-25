package commands

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/SeungheonOh/PixelOnTerminal/pixonterm"
	"github.com/SeungheonOh/PixelOnTerminal/screen"
	cv "gocv.io/x/gocv"
)

func init() {
	CommandMap["screen"] = &screenCommand{}
}

type screenCommand struct {
	fullScreen bool
	renderer   string
}

func (command *screenCommand) Description() string {
	return " [{X-Cord} {Y-Cord} {Width} {Height}] [Options]\n    Capture from screen with specified dimension\n    (currently only Xorg api supported)"
}

func (command *screenCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Print(USAGE, "\n  pix screen [{X-Cord} {Y-Cord} {Width} {Height}] [Options]\n\n")
		fmt.Print(OPTIONS, "\n")
		fs.PrintDefaults()
		fmt.Print("\n", DESCRIPTION, "\n  Capture from screen with specified dimension\n  (currently only Xorg api supported)\n")
		fmt.Println()

		os.Exit(1)
	}

	fs.BoolVar(&command.fullScreen, "f", false, "Make it fullscreen by stratching image")
	fs.StringVar(&command.renderer, "renderer", DEFAULT_RENDERER, "Select Renderer")

	return fs
}

func (command *screenCommand) Run(args []string) error {
	fs := command.FlagSet()
	fs.Parse(args[4:])

	if len(args) < 4 {
		fs.Usage()
		return nil
	}
	captureDimension := [4]int{}
	for i := 0; i < 4; i++ {
		number, err := strconv.Atoi(args[i])
		if err != nil {
			fs.Usage()
			return nil
		}
		captureDimension[i] = number
	}

	graber := screen.NewScreenGraber()

	running := true
	terminalSize := pixonterm.TermSize()
	go pixonterm.EventHandler(&running, &terminalSize)

	pixonterm.SetTerm()
	defer pixonterm.RecoverTerm()

	for running {
		capture, err := graber.Grab(captureDimension[0],
			captureDimension[1],
			captureDimension[2],
			captureDimension[3])
		if err != nil {
			return errors.New("Failed to grab screen")
		}

		img, err := cv.NewMatFromBytes(capture.Height, capture.Width, cv.MatTypeCV8UC3, capture.ToRGB())
		if err != nil {
			return errors.New("Failed to parse bytes")
		}
		defer img.Close()
		var imgSize image.Point
		if command.fullScreen {
			imgSize = image.Point{X: terminalSize.X, Y: terminalSize.Y * 2}
		} else {
			imgSize = pixonterm.CalculateSize(img, terminalSize)
		}
		cv.Resize(img, &img, imgSize, 0, 0, 1)

		err = pixonterm.PrintMat(img, command.renderer)
		if err != nil {
			return errors.New("failed to print image")
		}
	}

	return nil
}
