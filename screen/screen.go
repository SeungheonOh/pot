package screen

type ScreenCapture struct {
	Raw    []byte
	Width  int
	Height int
	X      int
	Y      int
}

func (S *ScreenCapture) ToRGB() []byte {
	ret := make([]byte, S.Width*S.Height*3)
	for i := 0; i < len(ret); i++ {
		ret[i] = S.Raw[i+(i/3)]
	}
	return ret
}
