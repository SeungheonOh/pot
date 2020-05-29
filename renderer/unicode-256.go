package renderer

import (
	"bytes"
	"fmt"
	"image"
)

func init() {
	RendererMap["unicode-256"] = NewUnicode256
}

type Unicode256 struct {
}

func NewUnicode256(options ...string) RenderEngine {
	return &Unicode256{}
}

func (r *Unicode256) Size(size image.Point) image.Point {
	return image.Point{size.X, size.Y * 2}
}

func (r *Unicode256) Render(img image.Image) (string, error) {
	var ret bytes.Buffer
	for i := 0; i < img.Bounds().Max.Y-2; i += 2 {
		for j := 0; j < img.Bounds().Max.X; j++ {
			r1, g1, b1, _ := img.At(j, i).RGBA()
			r2, g2, b2, _ := img.At(j, i+1).RGBA()
			fmt.Fprintf(&ret, "\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄\033[49m\033[39m",
				r1/257,
				g1/257,
				b1/257,
				r2/257,
				g2/257,
				b2/257)
		}
		fmt.Fprintf(&ret, "\n")
	}
	return ret.String(), nil
}

/*
func Unicode256(img cv.Mat, size image.Point) (string, error) {
	var buffer bytes.Buffer
	cv.Resize(img, &img, size, 0, 0, 1)

	imgPtr := img.DataPtrUint8()
	if img.Cols()*img.Rows()*3 != len(imgPtr) {
		return "", errors.New("Color RGB image is only supported")
	}

	for i := 0; i < img.Rows(); i += 2 {
		for j := 0; j < img.Cols()*3; j += 3 {
			if i+1 >= img.Rows() { // if height is not even, prevent to check pixel that isn't exist
				fmt.Fprintf(&buffer, "\033[48;2;%d;%d;%dm\033[30m▄\033[49m\033[39m",
					// Top Pixels
					imgPtr[i*img.Cols()*3+j+2],
					imgPtr[i*img.Cols()*3+j+1],
					imgPtr[i*img.Cols()*3+j])
			} else {
				fmt.Fprintf(&buffer, "\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄\033[49m\033[39m",
					// Top Pixels
					imgPtr[i*img.Cols()*3+j+2],
					imgPtr[i*img.Cols()*3+j+1],
					imgPtr[i*img.Cols()*3+j],
					// Bottom Pixels
					imgPtr[(i+1)*img.Cols()*3+j+2],
					imgPtr[(i+1)*img.Cols()*3+j+1],
					imgPtr[(i+1)*img.Cols()*3+j])
			}
		}
		// Don't draw a newline at the bottom of the terminal
		// Also clear all characters to the edge of the terminal
		// Prevents artifacts when the image.width < tty.width
		if i != img.Rows()-2 {
			fmt.Fprintf(&buffer, "\033[K\n")
		}
	}

	return buffer.String(), nil
}
*/
