package renderer

import (
	"fmt"
	"image"
	"time"
	//cv "gocv.io/x/gocv"
)

var RendererMap = make(map[string]func(...string) RenderEngine)

func init() {
	RendererMap["dummy"] = NewDummyRenderer
}

type RenderEngine interface {
	// Render
	Render(img image.Image) (string, error)
	// Returns image size that is required to make terminalsize of image
	Size(tsize image.Point) image.Point
}

type DummyRenderer struct {
}

func NewDummyRenderer(options ...string) RenderEngine {
	return &DummyRenderer{}
}

func (r *DummyRenderer) Size(tsize image.Point) image.Point {
	return image.Point{1, 1}
}

func (r *DummyRenderer) Render(img image.Image) (string, error) {
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintln("Dummy Renderer", "\nImage size: ", img.Bounds()), nil
}

/*
func Render(img cv.Mat, size image.Point, name string) error {
	return nil
		renderer, exist := RendererMap[name]
		if !exist {
			return errors.New("Renderer Not Found")
		}
		out, err := renderer(img, size)
		fmt.Fprintf(os.Stdout, "\033[0;0H")
		fmt.Fprintf(os.Stdout, out)
		// Clear from the end of the picture to the bottom of the tty
		// Also avoids leftover artifacts when image doesn't fill the tty
		fmt.Fprintf(os.Stdout, "\033[J")
		return err
}
*/
