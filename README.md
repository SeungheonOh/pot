# PixelOnTerminal
**These are Pixels, but in Terminals**
Light terminal media viewer using nothing but OpenCV

## Requirement
Download with `go get https://github.com/SeungheonOh/PixelOnTerminal.git`
or `git clone https://github.com/SeungheonOh/PixelOnTerminal.git`

build with ```go build main.go``` and copy/rename compile executable ```main``` to your PATH.

To build the project, you need 
```
Gocv
Crypto
Xlib(required only for screen subcommand)
```

## Renderers
### Why
Pix has very easy to add/configure renderers, fully modularized by files. User can write their own why to print Image to the terminal, like ascii-only renderer or no-unicode renderer.

### How to write one?
```renderer/ascii.go``` would be a great example. 
For additional information, it's getting ```cv.Mat```, which is already been resized for terminal size.
```cv.Mat``` has size of (Terminal Cols)x(Terminal Rows * 2), Rows are as twice as big as terminal size since one character in terminal takes of 
1x2 space(█) instead of 1x1. 

Render function can be used in many possible ways, including, but not limited to, process given image to print on the screen.
For example, Pixel on Terminal can be also used in X screen capturing with recording processor(or renderer).

## Subcommands
```
SUBCOMMANDS
	video   {File} [-Options]
		Play file on the terminal
		GIF format also does with this subcommand

	cam
		Cam subcommand loads webcam stream, print via Pixel On Terminal

	help    {SubCommand}
		Help Messages

	image   {File} [Options]
		Print file on the terminal
		subcommand url is equivlent to subcommand image with -u flag

	screen   [{X-Cord} {Y-Cord} {Width} {Height}] [Options]
		Capture from screen with specified dimension
		(currently only Xorg api supported)

	url   {URL} [Option]
		Print image from URL, equivlent of 'image' with '-d' option
```

## Examples

![WhiteOut](https://github.com/SeungheonOh/PixelOnTerminal/blob/master/doc/whiteout_small.png)
![Vim](https://github.com/SeungheonOh/PixelOnTerminal/blob/master/doc/vim_small.png)
![Youtube](https://github.com/SeungheonOh/PixelOnTerminal/blob/master/doc/youtube_small.png)

## Things to do (Maybe you can contribute)
- [x] Proper CLI system. 
- [x] Unified input selection for video, image, and cam.
- [ ] Get Some Stars (Yes press that Star NOW)

## Demo
[Demo at Reddit](https://www.reddit.com/r/unixporn/comments/d1gksi/oc_fully_terminal_based_webcamvideoimage_viewer/?utm_source=share&utm_medium=web2x)

## Inspired by
[Pastel](https://github.com/sharkdp/pastel) -> Definatly one of my favorate CLI app!

