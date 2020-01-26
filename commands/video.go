package commands

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/SeungheonOh/PixelOnTerminal/pixonterm"
	cv "gocv.io/x/gocv"
)

func init() {
	CommandMap["video"] = &videoCommand{}
}

type videoCommand struct {
	fullScreen bool
	repeat     bool
	loadUrl    bool
	renderer   string
}

func (command *videoCommand) Description() string {
	return " {File} [-Options]\n    Play file on the terminal\n    GIF format also does with this subcommand"
}

func (command *videoCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Print(USAGE, "\n  pix video {File} [Options]\n\n")
		fmt.Print(OPTIONS, "\n")
		fs.PrintDefaults()
		fmt.Print("\n", DESCRIPTION, "\n  Play file on the terminal\n  GIF format also does with this subcommand\n")
		fmt.Println()

		os.Exit(1)
	}

	fs.BoolVar(&command.fullScreen, "f", false, "Make it fullscreen by stratching image")
	fs.BoolVar(&command.repeat, "r", false, "Repeat the input until keyboard interupt")
	fs.StringVar(&command.renderer, "renderer", DEFAULT_RENDERER, "Select Renderer")

	return fs
}

func (command *videoCommand) Run(args []string) error {
	fs := command.FlagSet()

	if len(args) == 0 || args[0][0] == '-' {
		fs.Usage()
		return nil
	}
	file := args[0]
	fs.Parse(args[1:])

	stream, err := pixonterm.VideoStream(file)
	if err != nil || !stream.IsOpened() {
		return errors.New("failed to load video")
	}

	running := true
	frameRate := int(stream.Get(5))
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
		start := time.Now()

		img := cv.NewMat()
		defer img.Close()
		ok := stream.Read(&img)
		if !ok {
			if command.repeat {
				stream.Set(0, 0)
				continue
			}
			running = false
			continue
		}

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
		end := time.Now()
		time.Sleep(time.Second/time.Duration(frameRate) - end.Sub(start))
	}

	fmt.Print("\033c")
	return nil
}
