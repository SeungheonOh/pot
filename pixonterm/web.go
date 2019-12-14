package pixonterm

import (
	"errors"
	"io/ioutil"
	"net/http"

	cv "gocv.io/x/gocv"
)

func LoadFromURL(url string) (cv.Mat, error) {
	response, err := http.Get(url)
	if err != nil {
		return cv.NewMat(), errors.New("Failed load image")
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	img, err := cv.IMDecode(body, 1)

	if err != nil {
		return cv.NewMat(), errors.New("Failed decode image")
	}

	return img, nil
}
