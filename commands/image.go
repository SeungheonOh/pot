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
	CommandMap["image"] = &imageCommand{}
}

type imageCommand struct {
	fullScreen  bool
	repeatImage bool
	loadUrl     bool
}

func (command *imageCommand) Description() string {
	return " {File} [Options]\n    Print file on the terminal\n    subcommand url is equivlent to subcommand image with -u flag"
}

func (command *imageCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Print(USAGE, "\n  pix image {File} [Options]\n  pix url [URL] [Options]\n\n")
		fmt.Print(OPTIONS, "\n")
		fs.PrintDefaults()
		fmt.Print("\n", DESCRIPTION, "\n  Print file on the terminal\n  subcommand url is equivlent to subcommand image with -u flag\n")
		fmt.Println()

		os.Exit(1)
	}

	fs.BoolVar(&command.fullScreen, "f", false, "Make it fullscreen by stratching image")
	fs.BoolVar(&command.repeatImage, "r", false, "Repeat the input until keyboard interupt")
	fs.BoolVar(&command.loadUrl, "u", false, "Load from URL")

	return fs
}

func (command *imageCommand) Run(args []string) error {
	fs := command.FlagSet()

	if len(args) == 0 || args[0][0] == '-' {
		fs.Usage()
		return nil
	}
	file := args[0]
	fs.Parse(args[1:])

	running := true
	terminalSize := pixonterm.TermSize()
	if command.repeatImage {
		go pixonterm.EventHandler(&running, &terminalSize)
	}

	img := cv.NewMat()
	defer img.Close()
	if !command.loadUrl {
		img = cv.IMRead(file, cv.IMReadColor)
	} else {
		webImg, err := pixonterm.LoadFromURL(file)
		img = webImg
		if err != nil {
			return errors.New("failed to load image")
		}
	}
	if img.Empty() {
		return errors.New("failed to load image")
	}

	pixonterm.SetTerm()
	defer pixonterm.RecoverTerm()

	for run := true; run; run = command.repeatImage && running {
		var imgSize image.Point
		if command.fullScreen {
			imgSize = image.Point{X: terminalSize.X, Y: terminalSize.Y * 2}
		} else {
			imgSize = pixonterm.CalculateSize(img, terminalSize)
		}
		cv.Resize(img, &img, imgSize, 0, 0, 1)
		err := pixonterm.PrintMat(img)
		if err != nil {
			return errors.New("failed to print image")
		}
	}

	return nil
}
