package loader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func init() {
	LoaderMap["FFMPEG"] = NewFFMPEG
}

const (
	DEFAULT_FFMPEG_COMMAND_STRING = "ffmpeg -hide_banner -i %s -s %s -vf fps=15/1 -f rawvideo -c:v rawvideo -pix_fmt rgb24 -"
	// Dryrun is for getting total frames
	FFMPEG_DRYRUN = "ffmpeg -hide_banner -i %s -s 1x1 -vf fps=15/1 -f rawvideo -c:v rawvideo -pix_fmt rgb24 -"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0)
)

type FFMPEGJob struct {
	buf []byte
	img *image.Image
}

type FFMPEGWorkerPool struct {
	stop chan struct{}
	wg   sync.WaitGroup

	jobPool    chan *FFMPEGJob
	jobRequest chan *FFMPEGJob

	jobFn   func(*FFMPEGJob)
	bufsize int
}

func NewFFMPEGWorkerPool(bufsize int, poolsize int, jobFn func(*FFMPEGJob)) *FFMPEGWorkerPool {
	wp := &FFMPEGWorkerPool{
		stop:       make(chan struct{}),
		jobPool:    make(chan *FFMPEGJob, poolsize),
		jobRequest: make(chan *FFMPEGJob, poolsize),
		jobFn:      jobFn,
		bufsize:    bufsize,
	}

	for i := 0; i < poolsize; i++ {
		wp.jobPool <- &FFMPEGJob{
			buf: make([]byte, bufsize),
		}
	}

	wp.wg.Add(poolsize)

	for i := 0; i < poolsize; i++ {
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobRequest {
				wp.jobFn(job)
				wp.jobPool <- job
			}
		}()
	}

	return wp
}

func (wp *FFMPEGWorkerPool) Acquire() (*FFMPEGJob, chan<- *FFMPEGJob) {
	return <-wp.jobPool, wp.jobRequest
}

func (wp *FFMPEGWorkerPool) Wait() {
	close(wp.jobRequest)
	wp.wg.Wait()
}

type FFMPEG struct {
	commandstring string
}

func NewFFMPEG(options ...string) MediaLoader {
	if len(options) > 0 && options[0] != "" {
		return &FFMPEG{
			commandstring: options[0],
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

	cmds := strings.Split(fmt.Sprintf(l.commandstring, filename, fmt.Sprintf("%dx%d", size.X, size.Y)), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmderr, cmdout := &bytes.Buffer{}, &bytes.Buffer{}

	cmd.Stderr = cmderr
	cmd.Stdout = cmdout

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error when running ffmpeg binary: \n\t", cmderr.String()))
	}

	framesize, err := l.GetTotalFrames(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to fetch Frames: ", err))
	}
	frames := make([]image.Image, framesize)

	r := bufio.NewReader(cmdout)

	pool := NewFFMPEGWorkerPool(size.X*size.Y*3, runtime.NumCPU(), func(job *FFMPEGJob) {
		imgret := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{size.X, size.Y}})
		for i := 0; i < size.Y; i++ {
			for j := 0; j < size.X; j++ {
				imgret.Set(j, i, color.RGBA{
					job.buf[i*size.X*3+j*3],   // R
					job.buf[i*size.X*3+j*3+1], // G
					job.buf[i*size.X*3+j*3+2], // B
					255,                       // A
				})
			}
		}
		*job.img = imgret
	})

	for i := 0; i < framesize; i++ {
		job, ch := pool.Acquire()
		job.img = &frames[i]

		n, err := io.ReadFull(r, job.buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprint("Error on output reader: ", err))
		}

		ch <- job
	}
	pool.Wait()

	return frames, nil
}
