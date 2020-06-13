package loader

import (
	"fmt"
	"image"
	"testing"
	"time"
)

const (
	// This specific file has to have 901 frames when extracted with 30 frames per second
	//MEDIA_VIDEO = "https://file-examples.com/wp-content/uploads/2017/04/file_example_MP4_480_1_5MG.mp4"
	MEDIA_IMAGE = "https://raw.githubusercontent.com/mediaelement/mediaelement-files/master/big_buck_bunny.jpg"
	MEDIA_VIDEO = "./earth.mp4"
)

/*
func TestLoadingVideo(t *testing.T) {
	fmt.Println("==loading video==")
	loader := NewFFMPEG()
	_, err := loader.Load(MEDIA_VIDEO, image.Point{200, 100})
	if err != nil {
		t.Fatal(err)
		return
	}
}
*/

func TestVideo(t *testing.T) {
	start := time.Now()
	loader := NewFFMPEG()
	img_size := image.Point{500, 250}
	frames, err := loader.Load(MEDIA_VIDEO, img_size)
	if err != nil {
		t.Fatal(err)
		return
	}

	dur := time.Now().Sub(start)
	fmt.Println("loaded", len(frames), "frames in", dur, "at", img_size, "resolution", " fps: ", float64(len(frames)*1000)/float64(dur.Milliseconds()))
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
