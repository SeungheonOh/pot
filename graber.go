package main

/*
#cgo LDFLAGS: -L/usr/X11/lib -lX11
#cgo CFLAGS: -w
#include <string.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>

typedef struct {
	int X;
	int Y;
	int W;
	int H;
} WindowDimension;

WindowDimension GetDimensionById(Display *d, Window s, Window retval) {
	int x, y, h, w;
	Window child;
	XWindowAttributes xwa;
	XTranslateCoordinates( d, retval, XRootWindow(d, s), 0, 0, &x, &y, &child );
	XGetWindowAttributes( d, retval, &xwa );
	x = x - xwa.x;
	y = y - xwa.y;
	w = xwa.width;
	h = xwa.height;

	return (WindowDimension){ x, y, w, h };
}

XImage* GrabById(Display *d, Window s, Window retval){
	int x, y, h, w;
	Window child;
	XWindowAttributes xwa;
	XTranslateCoordinates( d, retval, XRootWindow(d, s), 0, 0, &x, &y, &child );
	XGetWindowAttributes( d, retval, &xwa );
	x = x - xwa.x;
	y = y - xwa.y;
	w = xwa.width;
	h = xwa.height;

	return XGetImage(d, XRootWindow(d, s), x, y, w, h, AllPlanes, ZPixmap);
}

XImage* Grab(Display *d, Window s, WindowDimension Win){
	return XGetImage(d, XRootWindow(d, s), Win.X, Win.Y, Win.W, Win.H, AllPlanes, ZPixmap);
}

Window GetRoot(Display *d){
	return DefaultScreen(d);
}

*/
import "C"

/*
XImage* Grab(Display *d, Window s, int x, int y, int w, int h){
	return XGetImage(d, XRootWindow(d, s), x, y, w, h, AllPlanes, ZPixmap);
}
*/

import (
	"unsafe"
)

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

type XScreenGraber struct {
	Display *C.Display
	Screen  C.Window
}

func NewXScreenGraber() *XScreenGraber {
	Display := C.XOpenDisplay(nil)
	Screen := C.GetRoot(Display)
	return &XScreenGraber{
		Display: Display,
		Screen:  Screen,
	}
}

func (X *XScreenGraber) Grab(x, y, w, h int) ScreenCapture {
	WinDimension := C.WindowDimension{C.int(w), C.int(y), C.int(w), C.int(h)}
	ximage := C.Grab(X.Display, X.Screen, WinDimension)
	pixels := C.GoBytes(unsafe.Pointer((*ximage).data), (*ximage).width*(*ximage).height*C.int(4))
	defer C.free(unsafe.Pointer((*ximage).data))
	defer C.free(unsafe.Pointer(ximage))

	return ScreenCapture{
		Raw:    pixels,
		Width:  int((*ximage).width),
		Height: int((*ximage).height),
		X:      x,
		Y:      y,
	}
}

func (X *XScreenGraber) GrabById(WinId uint64) ScreenCapture {
	WinDimension := C.GetDimensionById(X.Display, X.Screen, C.ulong(WinId))
	ximage := C.Grab(X.Display, X.Screen, WinDimension)
	pixels := C.GoBytes(unsafe.Pointer((*ximage).data), (*ximage).width*(*ximage).height*C.int(4))
	defer C.free(unsafe.Pointer((*ximage).data))
	defer C.free(unsafe.Pointer(ximage))

	return ScreenCapture{
		Raw:    pixels,
		Width:  int((*ximage).width),
		Height: int((*ximage).height),
		X:      int(WinDimension.X),
		Y:      int(WinDimension.Y),
	}
}

/*
func (X *XScreenGraber) Grab(WinId uint64) ([]byte, int, int) {
	ximage := C.Grab(X.Display, X.Screen, C.ulong(WinId))
	pixels := C.GoBytes(unsafe.Pointer((*ximage).data), (*ximage).width*(*ximage).height*C.int(4))
	defer C.free(unsafe.Pointer((*ximage).data))
	defer C.free(unsafe.Pointer(ximage))

	return pixels, int((*ximage).width), int((*ximage).height)
}
*/

func (X *XScreenGraber) Close() {
	C.XCloseDisplay(X.Display)
}

/*

func RGBAtoRGB(w, h int, raw []byte) ([]byte, error) {
	if w*h*4 != len(raw) {
		return nil, errors.New("Size and Data does not match")
	}
	ret := make([]byte, w*h*3)
	for i := 0; i < len(ret); i++ {
		ret[i] = raw[i+(i/3)]
	}
	return ret, nil
}
*/

/*
func main() {
	Graber := NewXScreenGraber()

	win := cv.NewWindow("cv")
	defer win.Close()

	WinId, err := strconv.ParseUint(os.Args[1], 16, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(WinId)

	for {
		arr, w, h := Graber.Grab(WinId)

		_ = arr
		fmt.Println(w, h)

		img, _ := cv.NewMatFromBytes(w, h, cv.MatTypeCV8UC4, arr)
		win.IMShow(img)
		win.WaitKey(1)
		img.Close()
	}
}

func main() {
	Display := C.XOpenDisplay(nil)
	Screen := C.GetRoot(Display)

	win := cv.NewWindow("cv")
	defer win.Close()

	WinId, err := strconv.ParseUint(os.Args[1], 16, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(WinId)

	for {
		ximage := C.Grab(Display, Screen, C.ulong(WinId))
		arr := C.GoBytes(unsafe.Pointer((*ximage).data), (*ximage).width*(*ximage).height*C.int(4))

		img, _ := cv.NewMatFromBytes(int((*ximage).height), int((*ximage).width), cv.MatTypeCV8UC4, arr)
		win.IMShow(img)
		win.WaitKey(1)
		img.Close()
		C.free(unsafe.Pointer((*ximage).data))
		C.free(unsafe.Pointer(ximage))
	}
}
*/
