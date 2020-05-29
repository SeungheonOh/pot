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

	"github.com/SeungheonOh/PixelOnTerminal/loader"
	"github.com/SeungheonOh/PixelOnTerminal/renderer"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	DEFAULT_RENDERER = ""
	DEFAULT_LOADER   = "FFMPEG"
)

type Context struct {
	filename string
	repeat   bool
	loader   loader.MediaLoader
	renderer renderer.RenderEngine
	size     image.Point

	buffer []image.Image
}

func NewContext(filename string, options FlagOptions) *Context {
	// Set default
	ctx := Context{
		filename: filename,
		repeat:   options.repeat,
		loader:   loader.LoaderMap[DEFAULT_LOADER](options.loaderoption),
		renderer: renderer.RendererMap["4x8"](options.rendereroption),
		size:     TermSize(),
	}

	if s := strings.Split(options.imagesize, "x"); len(s) == 2 {
		x, errx := strconv.Atoi(s[0])
		y, erry := strconv.Atoi(s[1])
		if errx == nil && erry == nil {
			ctx.size = image.Point{x, y}
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

func (c *Context) Load() error {
	frames, err := c.loader.Load(c.filename, c.renderer.Size(c.size))
	if err != nil {
		return err
	}
	c.buffer = frames
	return nil
}

type FlagOptions struct {
	renderer       string
	rendereroption string
	medialoader    string
	loaderoption   string
	imagesize      string
	repeat         bool
}

func FlagSet() (*flag.FlagSet, *FlagOptions) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Printf(`
PixelOnTerminal
	Usage: %s [FILE] [OPTIONS]...
	Print pixels in terminal screen
		`, filepath.Base(os.Args[0]))
		os.Exit(0)
	}

	option := &FlagOptions{}

	fs.StringVar(&option.renderer, "renderer", "", "Render engine")
	fs.StringVar(&option.rendereroption, "rendereroption", "", "Media loader option")
	fs.StringVar(&option.medialoader, "loader", "", "Media loader")
	fs.StringVar(&option.loaderoption, "loaderoption", "", "Media loader option")
	fs.StringVar(&option.imagesize, "size", "", "Image size")
	fs.BoolVar(&option.repeat, "r", false, "Repeat")

	return fs, option
}

func TermSize() image.Point {
	width, height, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return image.Point{width, height}
}

func main() {
	fs, options := FlagSet()
	if len(os.Args) < 2 {
		fs.Usage()
	}
	fs.Parse(os.Args[2:])

	fmt.Println(options)

	ctx := NewContext(os.Args[1], *options)
	fmt.Println(ctx)

	fmt.Fprintf(os.Stdout, "\033[2J\033[0;0HLoading frames to buffers")
	//defer fmt.Fprintf(os.Stdout, "\033[2J\033[?47l\0338")

	err := ctx.Load()
	if err != nil {
		panic(err)
	}

	var frame = make(chan image.Image)
	var refresh = time.NewTicker(time.Second / 15)
	defer refresh.Stop()

	go func() {
		for i := 0; i < len(ctx.buffer); i++ {
			if i+1 == len(ctx.buffer)-1 && ctx.repeat {
				i = 0
			}
			select {
			case frame <- ctx.buffer[i]:
				<-refresh.C
			case <-refresh.C:
			}

		}

		frame <- nil
	}()

	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGWINCH)

	for {
		//stop := timer()
		// handle signal
		var buf image.Image
		select {
		case sig := <-signals:
			switch sig {
			case os.Interrupt:
				//fmt.Fprintf(os.Stdout, "\033[2J\033[?47l\0338")
				close(signals)
				os.Exit(0)
			case syscall.SIGWINCH:
				//fmt.Fprintf(os.Stdout, "\033[2J\033[0;0HLoading")
				ctx.size = TermSize()
				err := ctx.Load()
				if err != nil {
					close(signals)
					os.Exit(0)
				}
			}
			continue
		case buf = <-frame:
		}

		//buf := ctx.buffer[i]
		if buf == nil {
			break
		}
		str, err := ctx.renderer.Render(buf)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\033[0;0H")
		fmt.Fprintf(os.Stdout, str)

		/*
			timeleft := time.Second/30 - stop()
			if timeleft >= 0 {
				time.Sleep(timeleft)
			} else {
				i += int(timeleft.Seconds() * -1 / (time.Second / 30).Seconds())
			}
		*/
	}
}
