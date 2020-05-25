package loader

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sync/semaphore"
)

const (
	DEFAULT_FFMPEG_COMMAND_STRING = "ffmpeg -i %s -s %s -vf fps=30/1 -f rawvideo -c:v rawvideo -pix_fmt rgb24 -"
	// Dryrun is for getting total frames
	FFMPEG_DRYRUN = "ffmpeg -i %s -s 1x1 -vf fps=30/1 -f rawvideo -c:v rawvideo -pix_fmt rgb24 -"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
)

type FFMPEG struct {
	commandstring string
}

func NewFFMPEG(option ...string) *FFMPEG {
	if len(option) > 0 {
		return &FFMPEG{
			commandstring: option[0],
		}
	}
	return &FFMPEG{
		commandstring: DEFAULT_FFMPEG_COMMAND_STRING,
	}
}

func (l *FFMPEG) GetTotalFrames(filename string) (int, error) {
	cmds := strings.Split(fmt.Sprintf(FFMPEG_DRYRUN, filename), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmderr, cmdout := &bytes.Buffer{}, &bytes.Buffer{}

	cmd.Stderr = cmderr
	cmd.Stdout = cmdout

	err := cmd.Run()
	if err != nil {
		return -1, errors.New(fmt.Sprint("Error when running ffmpeg binary: ", err, cmderr.String()))
	}

	framecount := 0
	r := bufio.NewReader(cmdout)
	buf := make([]byte, 3)
	for {
		n, err := io.ReadFull(r, buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}

		framecount++
	}

	return framecount, nil
}

/*
	Runs ffmpeg binary to ...
	if the file is video, get each frame of the video/gif
	if the file is image, get image
	with specified size.
	The frames are being loaded forehand and will get saved in the buffer.
	If the terminal size is altered, this whole process should happen again.
*/
func (l *FFMPEG) Load(filename string, size image.Point) ([]image.Image, error) {
	sem := semaphore.NewWeighted(int64(maxWorkers))

	cmds := strings.Split(fmt.Sprintf(l.commandstring, filename, fmt.Sprintf("%dx%d", size.X, size.Y)), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmderr, cmdout := &bytes.Buffer{}, &bytes.Buffer{}

	cmd.Stderr = cmderr
	cmd.Stdout = cmdout

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error when running ffmpeg binary: ", err, cmderr.String()))
	}

	framesize, err := l.GetTotalFrames(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to fetch Frames: ", err))
	}
	frames := make([]image.Image, framesize)

	indcounter := 0
	r := bufio.NewReader(cmdout)

	ctx := context.TODO()
	for {
		rawbuf := make([]byte, size.X*size.Y*3)
		n, err := io.ReadFull(r, rawbuf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprint("Error on output reader: ", err))
		}

		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, errors.New(fmt.Sprint("Concurrency Error: ", err))
		}

		go func(buf []byte, result *image.Image) {
			defer sem.Release(1)
			img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{size.X, size.Y}})
			for i := 0; i < size.Y; i++ {
				for j := 0; j < size.X; j++ {
					img.Set(j, i, color.RGBA{
						buf[i*size.X*3+j*3],   // R
						buf[i*size.X*3+j*3+1], // G
						buf[i*size.X*3+j*3+2], // B
						255,                   // A
					})
				}
			}
			result = img
		}(rawbuf, &frames[indcounter])
		indcounter++
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		return nil, errors.New(fmt.Sprint("Concurrency Error: ", err))
	}

	return frames, nil
}
