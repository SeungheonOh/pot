package loader

import (
	"image"
	"testing"
)

func Test(t *testing.T) {
	loader := NewGoimage()
	img_size := image.Point{500, 300}
	frames, err := loader.Load("./image1.jpg", img_size)
	if err != nil {
		t.Fatal(err)
		return
	}
	_ = frames
}
