package internal

import (
	"fmt"
	"os"
)

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
