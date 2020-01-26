package commands

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"os"

	"github.com/SeungheonOh/PixelOnTerminal/pixonterm"
	cv "gocv.io/x/gocv"
)

func init() {
	CommandMap["cam"] = &camCommand{}
}

type camCommand struct {
	fullScreen bool
	device     int
	renderer   string
}

func (command *camCommand) Description() string {
	return " \n    Cam subcommand loads webcam stream, print via Pixel On Terminal"
}

func (command *camCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Print(USAGE, "\n  pix cam [Options]\n\n")
		fmt.Print(OPTIONS, "\n")
		fs.PrintDefaults()
		fmt.Print("\n", DESCRIPTION, "\n  Cam subcommand loads webcam stream, print via Pixel On Terminal\n")
		fmt.Println()

		os.Exit(1)
	}

	fs.BoolVar(&command.fullScreen, "f", false, "Make it fullscreen by stratching image")
	fs.IntVar(&command.device, "d", 0, "Select Image capturing device")
	fs.StringVar(&command.renderer, "renderer", DEFAULT_RENDERER, "Select Renderer")

	return fs
}

func (command *camCommand) Run(args []string) error {
	fs := command.FlagSet()
	fs.Parse(args)

	webcam, err := pixonterm.WebCamStream(command.device)
	if err != nil {
		return errors.New("failed to open webcam")
	}

	running := true
	terminalSize := pixonterm.TermSize()
	go pixonterm.EventHandler(func() {
		running = false
		pixonterm.RecoverTerm()
	}, func() {
		terminalSize = pixonterm.TermSize()
	})

	pixonterm.SetTerm()
	defer pixonterm.RecoverTerm()

	for running {
		img := cv.NewMat()
		defer img.Close()
		webcam.Read(&img)
		var imgSize image.Point
		if command.fullScreen {
			imgSize = image.Point{X: terminalSize.X, Y: terminalSize.Y * 2}
		} else {
			imgSize = pixonterm.CalculateSize(img, terminalSize)
		}
		cv.Resize(img, &img, imgSize, 0, 0, 1)
		err := pixonterm.PrintMat(img, command.renderer)
		if err != nil {
			return err
		}
	}

	fmt.Print("\033c")
	return nil
}
