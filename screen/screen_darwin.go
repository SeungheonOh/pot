package screen

import (
	"errors"
)

/*
This have not been implemented yet
*/

type ScreenGraber struct {
}

func NewScreenGraber() *ScreenGraber {
	return &ScreenGraber{}
}

func (X *ScreenGraber) Grab(x, y, w, h int) (ScreenCapture, error) {
	return ScreenCapture{}, errors.New("Mac OS does now support screencapturing")
}

func (X *ScreenGraber) GrabById(WinId uint64) (ScreenCapture, error) {
	return ScreenCapture{}, errors.New("Mac OS does now support screencapturing")
}
