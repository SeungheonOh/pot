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

	"github.com/SeungheonOh/PixelOnTerminal/internal"
	"github.com/SeungheonOh/PixelOnTerminal/loader"
	"github.com/SeungheonOh/PixelOnTerminal/renderer"
)

const (
	DEFAULT_RENDERER = "4x8"
	DEFAULT_LOADER   = "FFMPEG"
)

type Context struct {
	filename string
	repeat   bool
	loader   loader.MediaLoader
	renderer renderer.RenderEngine
	size     image.Point

	buffer  []image.Image
	options FlagOptions
}

func NewContext(filename string, options FlagOptions) *Context {
	// Set default
	ctx := Context{
		filename: filename,
		repeat:   options.repeat,
		loader:   loader.LoaderMap[DEFAULT_LOADER](options.loaderoption),
		renderer: renderer.RendererMap[DEFAULT_RENDERER](options.rendereroption),

		options: options,
	}

	ctx.ReloadSize()

	if s := strings.Split(options.imagesize, "x"); len(s) == 2 {
		x, errx := strconv.Atoi(s[0])
		y, erry := strconv.Atoi(s[1])
		if errx == nil && erry == nil {
			ctx.size = image.Point{x, y}
		}
		if (x == -1 || y == -1) && y != x {
			size, err := ctx.loader.ImageSize(ctx.filename)
			if err == nil {
				ctx.size = internal.CalculateSizeWithRatio(ctx.size, size)
			}
		} else if y == x {
			ctx.options.imagesize = ""
			ctx.ReloadSize()
		}
	}

	if createloader, exist := loader.LoaderMap[options.medialoader]; exist {
		ctx.loader = createloader(options.loaderoption)
	}

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
	c.buffer = frames
	return nil
}

func (ctx *Context) Feeder(feed chan<- image.Image, fps int) {
	var refresh = time.NewTicker(time.Second / time.Duration(fps))
	defer refresh.Stop()

	for i := 0; i < len(ctx.buffer); i++ {
		if i+1 == len(ctx.buffer)-1 && ctx.repeat {
			i = 0
		}
		select {
		case feed <- ctx.buffer[i]:
			<-refresh.C
		case <-refresh.C:
		}
	}
	feed <- nil
}

func (ctx *Context) Loop() {
	var feed = make(chan image.Image)
	defer close(feed)

	go ctx.Feeder(feed, 15)

	var signals = make(chan os.Signal, 1)
	defer close(signals)
	signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)

	for {
		select {
		case sig := <-signals:
			switch sig {
			case os.Interrupt:
				return
			case syscall.SIGWINCH:
				if ctx.options.imagesize == "" {
					ctx.ReloadSize()
					err := ctx.Load()
					if err != nil {
						internal.ClearScreen()
						return
					}
				}
			}
			continue
		case buf := <-feed:
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
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

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
