package screen

/*
This have not been implemented yet
*/

type ScreenGraber struct {
}

func NewScreenGraber() *XScreenGraber {
	return &Dummy{}
}

func (X *ScreenGraber) Grab(x, y, w, h int) (ScreenCapture, error) {
	return nil, errors.New("Your OS does now support screencapturing")
}

func (X *ScreenGraber) GrabById(WinId uint64) (ScreenCapture, error) {
	return nil, errors.New("Your OS does now support screencapturing")
}
