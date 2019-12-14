# PixelOnTerminal
Light terminal media viewer using nothing but OpenCV

## Requirement
You need 24bit color supported terminal emulator, such as `st`, `xterm`, `terminator`... : )

Download with `go get https://github.com/SeungheonOh/PixelOnTerminal.git`
or `git clone https://github.com/SeungheonOh/PixelOnTerminal.git`

To build the project, you need 
```
Gocv
Crypto
Xlib(required only for screen subcommand)
```

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

