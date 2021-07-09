package renderer

import (
	"fmt"
	"image"
	"time"
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
