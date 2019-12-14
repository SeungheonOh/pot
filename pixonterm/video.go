package pixonterm

import (
	cv "gocv.io/x/gocv"
)

func WebCamStream(id int) (*cv.VideoCapture, error) {
	return cv.VideoCaptureDevice(id)
}

func VideoStream(file string) (*cv.VideoCapture, error) {
	return cv.OpenVideoCapture(file)
}
