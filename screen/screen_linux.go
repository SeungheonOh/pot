package screen

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

import (
	"unsafe"
)

type ScreenGraber struct {
	Display *C.Display
	Screen  C.Window
}

func NewScreenGraber() *ScreenGraber {
	Display := C.XOpenDisplay(nil)
	Screen := C.GetRoot(Display)
	return &ScreenGraber{
		Display: Display,
		Screen:  Screen,
	}
}

func (X *ScreenGraber) Grab(x, y, w, h int) (ScreenCapture, error) {
	WinDimension := C.WindowDimension{C.int(x), C.int(y), C.int(w), C.int(h)}
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
	}, nil
}

func (X *ScreenGraber) GrabById(WinId uint64) (ScreenCapture, error) {
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
	}, nil
}

func (X *ScreenGraber) Close() {
	C.XCloseDisplay(X.Display)
}
