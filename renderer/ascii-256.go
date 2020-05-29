package renderer

import (
	"bytes"
	"fmt"
	"image"
)

func init() {
	RendererMap["ascii-256"] = NewAscii256
}

type Ascii256 struct {
}

func NewAscii256(options ...string) RenderEngine {
	return &Ascii256{}
}

func (r *Ascii256) Size(size image.Point) image.Point {
	return image.Point{size.X, size.Y}
}

func (r *Ascii256) Render(img image.Image) (string, error) {
	var ret bytes.Buffer
	for i := 0; i < img.Bounds().Max.Y-1; i++ {
		for j := 0; j < img.Bounds().Max.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			fmt.Fprintf(&ret, "\033[38;2;%d;%d;%dm█\033[39m",
				r/257,
				g/257,
				b/257)
		}
		fmt.Fprintf(&ret, "\n")
	}
	return ret.String(), nil
}
