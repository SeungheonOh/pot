package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SeungheonOh/pot/internal"
	"github.com/SeungheonOh/pot/loader"
	"github.com/SeungheonOh/pot/renderer"
)

const (
	DEFAULT_RENDERER = "4x8"
	DEFAULT_LOADER   = "FFMPEG"
)

/*
	This struct stores all informations and modules that are required.
	It has to be kept as simple as possible to make the modular system clean.
*/
type Context struct {
	filename string
	loader   loader.MediaLoader
	renderer renderer.RenderEngine
	size     image.Point

	buffer  []image.Image
	options FlagOptions
}

/*
	Creating new Context according to FlagOptions.
*/
func NewContext(filename string, options FlagOptions) *Context {
	// Default settings
	ctx := Context{
		filename: filename,
		loader:   loader.LoaderMap[DEFAULT_LOADER](options.loaderoption),
		renderer: renderer.RendererMap[DEFAULT_RENDERER](options.rendereroption),

		options: options,
	}

	// Initiate Image size
	ctx.ReloadSize()

	// Custom size specified by --size
	if s := strings.Split(options.imagesize, "x"); len(s) == 2 {
		x, errx := strconv.Atoi(s[0])
		y, erry := strconv.Atoi(s[1])
		if errx == nil && erry == nil {
			ctx.size = image.Point{x, y}
		}

		// Handle when width or height given is -1, if it is -1 this will
		// automatically set size of other side according to the ratio of image
		if (x == -1 || y == -1) && y != x {
			size, err := ctx.loader.ImageSize(ctx.filename)
			if err == nil {
				ctx.size = internal.CalculateSizeWithRatio(ctx.size, size)
			}
		} else if y == x { // Handle when both width and height given is -1.
			ctx.options.imagesize = ""
			ctx.ReloadSize()
		}
	}

	// When option for non-default loader is provided, use it
	if createloader, exist := loader.LoaderMap[options.medialoader]; exist {
		ctx.loader = createloader(options.loaderoption)
	}

	// When option for non-default renderer is provided, use it
	if createrenderer, exist := renderer.RendererMap[options.renderer]; exist {
		ctx.renderer = createrenderer(options.rendereroption)
	}

	return &ctx
}

/*
	Reloads size stored in context. It will automatically calculate ratio between image and terminalize and
	return accordingly unless option -f is specified.
*/
func (c *Context) ReloadSize() {
	// Don't need to set size according to the ratio of original image when -f flag is on
	if !c.options.fullscreen {
		size, err := c.loader.ImageSize(c.filename)
		if err == nil {
			c.size = internal.CalculateSizeWithRatio(internal.TermSize(), size)
			return
		}
	}
	c.size = internal.TermSize()
}

/*
	Loads each frame of video/gif/image to the context.buffer
*/
func (c *Context) Load() error {
	frames, err := c.loader.Load(c.filename, c.renderer.Size(c.size))
	if err != nil {
		return err
	}
	c.buffer = nil
	c.buffer = frames
	return nil
}

/*
	Feeder feeds each frame to the loop itself according to specified refresh rate.
	The pace is kept by time.Ticker. It has to send nil at the end of the stream
	to indicate end of stream.
*/
func (ctx *Context) Feeder(feed chan<- image.Image, fps int) {
	// Pace maker
	var refresh = time.NewTicker(time.Second / time.Duration(fps))
	defer refresh.Stop()

	for i := 0; i < len(ctx.buffer); i++ {
		// If -r flag is on, loops back to the 0th frame at the end
		if i+1 == len(ctx.buffer)-1 && ctx.options.repeat {
			i = 0
		}
		select {
		case feed <- ctx.buffer[i]:
			<-refresh.C // wait for ticker
		case <-refresh.C: // If render engine takes more time than 1/fps second, it has to skip frame
		}
	}

	feed <- nil // end of stream
}

/*
	Render Loop renders frame fed from Feeder and take care of syscalls.
*/
func (ctx *Context) Loop() {

	// Channel initiation
	var feed = make(chan image.Image)
	var signals = make(chan os.Signal, 1)
	defer close(signals)
	defer close(feed)
	signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)
	go ctx.Feeder(feed, 15)

	for {
		select {
		// Syscalls
		case sig := <-signals:
			switch sig {
			case os.Interrupt: // Stop at interrupt
				return
			case syscall.SIGWINCH: // Reload at Sigwinch
				if ctx.options.imagesize == "" { // only if size is not specified
					ctx.ReloadSize()
					err := ctx.Load()
					if err != nil {
						internal.ClearScreen()
						return
					}
				}
			}
			continue
		// Frames
		case buf := <-feed:
			// stop at the end of stream
			if buf == nil {
				return
			}
			str, err := ctx.renderer.Render(buf)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}
			internal.SetCursorPos(0, 0)
			fmt.Fprintf(os.Stdout, str)
		}
	}
}

/*
	Stores options with FlagOptions struct
*/
type FlagOptions struct {
	renderer       string
	rendereroption string
	medialoader    string
	loaderoption   string
	imagesize      string
	repeat         bool
	fullscreen     bool
}

func FlagSet() (*flag.FlagSet, *FlagOptions) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Printf(`
PixelOnTerminal
Print pixels in terminal screen

Usage: 
  %s [FILE] [OPTIONS]...

Options:
`, filepath.Base(os.Args[0]))

		fs.PrintDefaults()
		os.Exit(0)
	}

	option := &FlagOptions{}

	fs.StringVar(&option.renderer, "renderer", "", "Render engine")
	fs.StringVar(&option.rendereroption, "rendereroption", "", "Media loader option")
	fs.StringVar(&option.medialoader, "loader", "", "Media loader")
	fs.StringVar(&option.loaderoption, "loaderoption", "", "Media loader option")
	fs.StringVar(&option.imagesize, "size", "", "Image size")
	fs.BoolVar(&option.repeat, "r", false, "Repeat")
	fs.BoolVar(&option.fullscreen, "f", false, "Repeat")

	return fs, option
}

func main() {
	fs, options := FlagSet()
	if len(os.Args) < 2 {
		fs.Usage()
	}
	fs.Parse(os.Args[2:])

	ctx := NewContext(os.Args[1], *options)

	err := ctx.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
		os.Exit(1)
	}

	// TTY setup
	internal.ClearScreen()
	internal.Cursor(false)
	internal.SetCursorPos(0, 0)
	defer internal.Cursor(true)

	ctx.Loop()
}

/*
4x8
1 0 1 0
0 1 0 1
1 0 1 0
0 1 0 1
1 0 1 0
0 1 0 1
1 0 1 0
0 1 0 1

2x4
1 1
1 1
1 1
1 1

4x8
1 1 0 0
1 1 0 0
1 1 0 0
1 1 0 0
0 0 1 1
0 0 1 1
0 0 1 1
0 0 1 1

2x4
1 0
1 0
0 1
0 1

1x2
1
1

*/
