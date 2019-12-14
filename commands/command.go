package commands

import (
	"errors"
	"flag"
)

const (
	USAGE       = "\u001b[32mUSAGE\u001b[0m"
	SUBCOMMANDS = "\u001b[32mSUBCOMMANDS\u001b[0m"
	OPTIONS     = "\u001b[32mOPTIONS\u001b[0m"
	DESCRIPTION = "\u001b[32mDESCRIPTION\u001b[0m"
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
		return errors.New("Unknown command")
	}

	return cmd.Run(subarg)
}
