package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func init() {
	CommandMap["help"] = &helpCommand{}
}

type helpCommand struct {
}

func (command *helpCommand) Description() string {
	return "  {SubCommand}\n    Help Messages"
}

func (command *helpCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Println(`┌─────────────────────────┐`)
		fmt.Println(`│┌─┐┬┐ ┌┌─┐┌┐┌┌┬┐┌─┐┬─┐┌┬┐│`)
		fmt.Println(`│├─┘│└─┐│ ││││ │ ├┤ ├┬┘││││`)
		fmt.Println(`│┴  ┴┘ └└─┘┘└┘ ┴ └─┘┴└─┴ ┴│`)
		fmt.Println(`└─────────────────────────┘`)
		fmt.Println("PixelOnTerminal - CLI media viewer")

		fmt.Print("\n")
		fmt.Print(USAGE, "\n  pix [SubCommand] [Arguments] [Options]\n\n")
		fmt.Print(SUBCOMMANDS, "\n")
		for cmdName := range CommandMap {
			fmt.Print("  \u001b[33m", cmdName, "\u001b[0m")
			cmd, _ := CommandMap[cmdName]
			fmt.Print("  ", cmd.Description(), "\n\n")
		}
		fs.PrintDefaults()
		fmt.Print("\nBy Seungheon Oh, 2019\nUnder MIT License\n")
		fmt.Println()

		os.Exit(0)
	}

	return fs
}

func (command *helpCommand) Run(args []string) error {
	fs := command.FlagSet()
	if len(args) == 0 {
		fs.Usage()
		return nil
	}

	help, exist := CommandMap[args[0]]
	if !exist {
		return errors.New("Unknown Command")
	}
	help.Run([]string{"-help"})

	return nil
}
