package loader

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

const (
	// This specific file has to have 901 frames when extracted with 30 frames per second
	//MEDIA_VIDEO = "https://file-examples.com/wp-content/uploads/2017/04/file_example_MP4_480_1_5MG.mp4"
	MEDIA_IMAGE = "https://raw.githubusercontent.com/mediaelement/mediaelement-files/master/big_buck_bunny.jpg"
	MEDIA_VIDEO = "./earth.mp4"
)

func TestLoadingVideo(t *testing.T) {
	fmt.Println("==loading video==")
	loader := NewFFMPEG()
	fmt.Println(loader.ImageSize(MEDIA_VIDEO))
	frames, err := loader.Load(MEDIA_VIDEO, image.Point{200, 100})
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(frames) != 901 {
		t.Fatal("Frame number doesn't match", len(frames))
	}
	return
	for i, img := range frames {
		f, _ := os.Create(fmt.Sprintf("%d%d%d.png", i/100, (i%100)/10, i%10))
		png.Encode(f, img)
	}
}

/*
func TestLoadingImage(t *testing.T) {
	fmt.Println("==loading image==")
	loader := NewFFMPEG()
	frames, err := loader.Load(MEDIA_IMAGE, image.Point{200, 100})
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(frames) != 1 {
		t.Fatal("Frame number doesn't match")
	}
}
*/
