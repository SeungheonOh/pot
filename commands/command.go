package commands

import (
	"flag"
)

const (
	USAGE            = "\u001b[32mUSAGE\u001b[0m"
	SUBCOMMANDS      = "\u001b[32mSUBCOMMANDS\u001b[0m"
	OPTIONS          = "\u001b[32mOPTIONS\u001b[0m"
	DESCRIPTION      = "\u001b[32mDESCRIPTION\u001b[0m"
	DEFAULT_RENDERER = "unicode-256"
)

var CommandMap = make(map[string]Command)

type Command interface {
	Description() string
	FlagSet() *flag.FlagSet
	Run(args []string) error
}

func Run(args []string) error {
	if len(args) <= 1 {
		args = append(args, "help")
	}
	subcommand := args[1]
	subarg := args[2:]

	cmd, exist := CommandMap[subcommand]
	if !exist {
		//return errors.New("Unknown command")
		cmd, err := DetermineSubcommand(args[1:])
		if err != nil {
			return err
		}
		return cmd.Run(args[1:])
	}

	return cmd.Run(subarg)
}
