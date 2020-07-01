package loader

import (
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"os"
)

func init() {
	LoaderMap["GoimageLoader"] = NewGoimage
	LoaderMap["std"] = NewGoimage
}

func Lerp(start, end, t float64) float64 {
	return start*(1-t) + end*t
}

type GoimageLoader struct {
}

func NewGoimage(options ...string) MediaLoader {
	return &GoimageLoader{}
}

func (l *GoimageLoader) ImageSize(filename string) (image.Point, error) {
	imgfile, err := os.Open(filename)
	if err != nil {
		return image.Point{-1, -1}, errors.New("Failed to fetch size")
	}
	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	if err != nil {
		return image.Point{-1, -1}, errors.New("Failed to decode size")
	}

	return img.Bounds().Max, nil
}

func (l *GoimageLoader) Load(filename string, size image.Point) ([]image.Image, error) {
	imgfile, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("Failed to open file")
	}
	defer imgfile.Close()

	img, _, err := image.Decode(imgfile)
	if err != nil {
		return nil, errors.New("Failed to decode image")
	}
	ret := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{size.X, size.Y}})

	for x := 0; x < ret.Bounds().Max.X; x++ {
		for y := 0; y < ret.Bounds().Max.Y; y++ {
			sx := int(math.Round((Lerp(0, float64(img.Bounds().Max.X), float64(x)/float64(ret.Bounds().Max.X)))))
			sy := int(math.Round((Lerp(0, float64(img.Bounds().Max.Y), float64(y)/float64(ret.Bounds().Max.Y)))))

			r, g, b, _ := img.At(sx, sy).RGBA()
			ret.Set(x, y, color.RGBA{
				uint8(r / 257),
				uint8(g / 257),
				uint8(b / 257),
				255,
			})
		}
	}

	return []image.Image{ret}, nil
}
