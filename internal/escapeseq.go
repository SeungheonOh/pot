package internal

import (
	"fmt"
	"image"
	"math"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func TermSize() image.Point {
	width, height, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return image.Point{width, height}
}

func CalculateSizeWithRatio(term, size image.Point) image.Point {
	if term.X == -1 && term.Y == -1 {
		return term
	}
	if term.X == -1 {
		term.X = math.MaxInt32
	}
	if term.Y == -1 {
		term.Y = math.MaxInt32
	}

	termRatio := float64(term.Y) / float64(term.X)
	imgRatio := float64(size.Y) / float64(size.X*2)
	if imgRatio > termRatio {
		return image.Point{int(float64(term.Y) / imgRatio), term.Y}
	} else {
		return image.Point{term.X, int(imgRatio * float64(term.X))}
	}
}

func Cursor(con bool) {
	if con {
		fmt.Fprintf(os.Stdout, "\033[?25h")
	} else {
		fmt.Fprintf(os.Stdout, "\033[?25l")
	}
}

func SetCursorPos(x, y uint) {
	fmt.Fprintf(os.Stdout, "\033[%d;%dH", x, y)
}

func ClearScreen() {
	fmt.Fprintf(os.Stdout, "\033[2J")
}
