package main

import (
	"fmt"
	"os"

	"github.com/SeungheonOh/PixelOnTerminal/commands"
)

func main() {
	err := commands.Run(os.Args)
	if err != nil {
		fmt.Println("\n\u001b[31mERROR: \u001b[0m", err)
		os.Exit(1)
	}
}
