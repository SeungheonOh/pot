package commands

import (
	"errors"
	"flag"
)

func init() {
	CommandMap["url"] = &urlCommand{}
}

type urlCommand struct {
}

func (command *urlCommand) Description() string {
	return " {URL} [Option]\n    Print image from URL, equivlent of 'image' with '-d' option"
}

func (command *urlCommand) FlagSet() *flag.FlagSet {
	return nil
}

func (command *urlCommand) Run(args []string) error {
	imageCmd, exist := CommandMap["image"]
	if !exist {
		return errors.New("Internal Error: Failed to find image command")
	}

	return imageCmd.Run(append(args, "-u"))
}
